package lr35902

import (
	"fmt"
	"testing"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/stretchr/testify/assert"
)

// New helper: setupWithOpcode initializes CPU and writes opcodes + extra bytes.
func setupWithOpcode(codes ...uint8) (*memory.Bus, *LR35902) {
	bus := &memory.Bus{}
	io := &registers.Registers{}
	bus.AddDevice(0x0000, 0xFFFF, &memory.Memory{Buffer: make([]byte, 0x10000)})
	cpu := NewLR35902(bus, io)
	// Write provided opcodes to PC sequentially.
	pc := cpu.registers.pc
	for i, code := range codes {
		bus.Write(pc+uint16(i), code)
	}
	return bus, cpu
}

func TestInstructions(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Run("Instruction: NOP", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x00)
			initPC := cpu.registers.pc
			cpu.Step()
			assert.Equal(t, initPC+1, cpu.registers.pc, "NOP should increment PC by 1")
		})

		t.Run("Instruction: LD_BC_d16", func(t *testing.T) {
			// Pass opcode and immediate 16-bit little-endian bytes.
			_, cpu := setupWithOpcode(0x01, 0x42, 0x24)
			cpu.Step()
			assert.Equal(t, uint16(0x2442), toRegisterPair(cpu.registers.b, cpu.registers.c), "BC should load immediate 16-bit value")
		})
	})

	t.Run("8-Bit Loads", func(t *testing.T) {
		t.Run("Instruction: LD_d8", func(t *testing.T) {
			{
				_, cpu := setupWithOpcode(0x06, 0x42)
				cpu.Step()
				assert.Equal(t, uint8(0x42), cpu.registers.b, "B should load immediate 8-bit value")
			}
			{
				_, cpu := setupWithOpcode(0x0E, 0x55)
				cpu.Step()
				assert.Equal(t, uint8(0x55), cpu.registers.c, "C should load immediate 8-bit value")
			}
		})
		// Example for an operation test:
		t.Run("Operation: INC B", func(t *testing.T) {
			{
				_, cpu := setupWithOpcode(0x04)
				cpu.registers.b = 1
				cpu.Step()
				assert.Equal(t, uint8(2), cpu.registers.b, "B should increment by 1")
				assert.False(t, cpu.flags.Zero, "Z flag should be reset on INC")
				assert.False(t, cpu.flags.Subtract, "N flag should be reset on INC")
				assert.False(t, cpu.flags.HalfCarry, "H flag should be reset on INC")
			}
		})
		t.Run("Operation: DEC B", func(t *testing.T) {
			decBTests := []struct {
				name      string
				initVal   uint8
				expResult uint8
				expZero   bool
			}{
				{"DEC_B_from_0", 0, 0xFF, false},
				{"DEC_B_from_1", 1, 0, true},
				{"DEC_B_from_2", 2, 1, false},
			}
			for _, tt := range decBTests {
				t.Run(fmt.Sprintf("DEC_B_%s", tt.name), func(t *testing.T) {
					_, cpu := setupWithOpcode(0x05)
					cpu.registers.b = tt.initVal
					cpu.Step()
					assert.Equal(t, tt.expResult, cpu.registers.b, "DEC B did not produce expected value")
					assert.Equal(t, tt.expZero, cpu.flags.Zero, "Zero flag mismatch on DEC B")
					assert.True(t, cpu.flags.Subtract, "N flag should be set on DEC")
				})
			}
		})
	})

	t.Run("16-Bit Operations", func(t *testing.T) {
		t.Run("Instruction: INC/DEC BC", func(t *testing.T) {
			{
				_, cpu := setupWithOpcode(0x03)
				cpu.registers.b, cpu.registers.c = fromRegisterPair(0x01)
				cpu.Step()
				assert.Equal(t, uint16(2), toRegisterPair(cpu.registers.b, cpu.registers.c), "BC should increment by 1")
			}
			decBCTests := []struct {
				name      string
				initBC    uint16
				expResult uint16
			}{
				{"DEC_BC_from_1", 0x0001, 0x0000},
				{"DEC_BC_from_2", 0x0002, 0x0001},
			}
			for _, tt := range decBCTests {
				t.Run(fmt.Sprintf("DEC_BC_%s", tt.name), func(t *testing.T) {
					_, cpu := setupWithOpcode(0x0B)
					cpu.registers.b, cpu.registers.c = fromRegisterPair(tt.initBC)
					cpu.Step()
					res := toRegisterPair(cpu.registers.b, cpu.registers.c)
					assert.Equal(t, tt.expResult, res, "DEC BC did not produce expected value")
				})
			}
		})
		t.Run("Operation: ADD HL,BC", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x09)
			cpu.registers.h, cpu.registers.l = fromRegisterPair(0x1)
			cpu.registers.b, cpu.registers.c = fromRegisterPair(0x1)
			cpu.Step()
			assert.Equal(t, uint16(0x2), toRegisterPair(cpu.registers.h, cpu.registers.l), "HL should add BC")
			assert.False(t, cpu.flags.Subtract, "N flag should be reset in addition")
		})
	})

	t.Run("Memory", func(t *testing.T) {
		t.Run("Instruction: Address Operations", func(t *testing.T) {
			t.Run("LD_A_from_BC", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0x02)
				cpu.registers.b, cpu.registers.c = fromRegisterPair(0x01)
				cpu.registers.a = 0xAA
				cpu.Step()
				assert.Equal(t, uint8(0xAA), bus.Read(0x0001), "Memory at address BC should be loaded with A")
			})
			t.Run("LD_A_from_a16", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0x0A)
				cpu.registers.b = 0x01
				cpu.registers.c = 0x00
				bus.Write(0x100, 0xBB)
				cpu.Step()
				assert.Equal(t, uint8(0xBB), cpu.registers.a, "A should load value from memory at address BC")
			})
			t.Run("LD_a16_from_A", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0x08, 0x0B, 0x00)
				cpu.registers.sp = 0x1234
				cpu.Step()
				assert.Equal(t, uint8(0x34), bus.Read(0x0B), "Memory low should be SP's low byte")
				assert.Equal(t, uint8(0x12), bus.Read(0x0C), "Memory high should be SP's high byte")
			})
		})
		t.Run("Instruction: LD Variants", func(t *testing.T) {
			t.Run("LD_A_from_a16", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xFA)
				addr := uint16(0x2000)
				bus.Write(addr, 0x7F)
				bus.Write(cpu.registers.pc+1, uint8(addr&0xFF))
				bus.Write(cpu.registers.pc+2, uint8(addr>>8))
				cpu.Step()
				assert.Equal(t, uint8(0x7F), cpu.registers.a, "LD A,(a16) should load value from memory")
			})
			t.Run("LD_a16_from_A", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xEA)
				addr := uint16(0x3000)
				cpu.registers.a = 0x3C
				bus.Write(cpu.registers.pc+1, uint8(addr&0xFF))
				bus.Write(cpu.registers.pc+2, uint8(addr>>8))
				cpu.Step()
				assert.Equal(t, uint8(0x3C), bus.Read(addr), "LD (a16),A should store A into memory")
			})
			t.Run("LDH_A_from_n", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xF0)
				offset := uint8(0x20)
				addr := uint16(0xFF00) + uint16(offset)
				bus.Write(addr, 0x99)
				bus.Write(cpu.registers.pc+1, offset)
				cpu.Step()
				assert.Equal(t, uint8(0x99), cpu.registers.a, "LDH A,(n) should load value from 0xFF00+n")
			})
			t.Run("LDH_n_from_A", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xE0)
				offset := uint8(0x30)
				cpu.registers.a = 0xAB
				addr := uint16(0xFF00) + uint16(offset)
				bus.Write(cpu.registers.pc+1, offset)
				cpu.Step()
				assert.Equal(t, uint8(0xAB), bus.Read(addr), "LDH (n),A should store A into memory at 0xFF00+n")
			})
		})
	})

	t.Run("Rotation", func(t *testing.T) {
		t.Run("Instruction: RLCA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x07)
			cpu.registers.a = 0x80
			cpu.Step()
			assert.Equal(t, uint8(0x01), cpu.registers.a, "A should rotate left (RLCA)")
			assert.True(t, cpu.flags.Carry, "Carry flag should be set by RLCA")
		})
		t.Run("Instruction: RRCA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x0F)
			cpu.registers.a = 0x01
			cpu.Step()
			assert.Equal(t, uint8(0x80), cpu.registers.a, "A should rotate right (RRCA)")
			assert.True(t, cpu.flags.Carry, "Carry flag should be set by RRCA")
		})
		t.Run("Instruction: RLA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x17)
			cpu.registers.a = 0x80
			cpu.flags.Carry = true
			cpu.Step()
			assert.Equal(t, uint8(0x01), cpu.registers.a, "A should rotate left (RLA)")
			assert.True(t, cpu.flags.Carry, "Carry flag should be set by RLA")
		})
		t.Run("Instruction: RRA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x1F)
			cpu.registers.a = 0x01
			cpu.flags.Carry = true
			cpu.Step()
			assert.Equal(t, uint8(0x80), cpu.registers.a, "A should rotate right (RRA)")
			// Updated: Carry flag should remain true after RRA.
			assert.True(t, cpu.flags.Carry, "Carry flag should be set by RRA")
		})
	})

	t.Run("ALU", func(t *testing.T) {
		t.Run("Instruction: Arithmetic & Logic", func(t *testing.T) {
			// ...existing table-driven ALU tests...
			type flags struct {
				zero, carry, halfCarry, subtract bool
			}
			type testCase struct {
				name            string
				opcode          uint8
				initA, initB    uint8
				initCarry       bool
				expectedA       uint8
				expectedFlags   flags
				checkAUnchanged bool
			}
			tests := []testCase{
				{"ADD_A_B_simple", 0x80, 1, 2, false, 3, flags{false, false, false, false}, false},
				{"ADD_A_B_overflow", 0x80, 0xFF, 1, false, 0, flags{true, true, true, false}, false},
				{"ADC_A_B_simple", 0x88, 1, 2, true, 4, flags{false, false, false, false}, false},
				{"ADC_A_B_overflow", 0x88, 0xFF, 0, true, 0, flags{true, true, true, false}, false},
				{"SUB_A_B_simple", 0x90, 5, 3, false, 2, flags{false, false, false, true}, false},
				{"SUB_A_B_zero", 0x90, 3, 3, false, 0, flags{true, false, false, true}, false},
				{"AND_A_B", 0xA0, 0x55, 0xF0, false, 0x50, flags{false, false, true, false}, false},
				{"XOR_A_B", 0xA8, 0xFF, 0x0F, false, 0xF0, flags{false, false, false, false}, false},
				{"OR_A_B", 0xB0, 0x55, 0xAA, false, 0xFF, flags{false, false, false, false}, false},
				{"CP_A_B_equal", 0xB8, 3, 3, false, 3, flags{true, false, false, true}, true},
				{"CP_A_B_diff", 0xB8, 4, 3, false, 4, flags{false, false, false, true}, true},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					_, cpu := setupWithOpcode(tc.opcode)
					cpu.registers.a = tc.initA
					cpu.registers.b = tc.initB
					cpu.flags.Carry = tc.initCarry
					cpu.Step()
					if !tc.checkAUnchanged {
						assert.Equal(t, tc.expectedA, cpu.registers.a, tc.name+": A value")
					} else {
						assert.Equal(t, tc.initA, cpu.registers.a, tc.name+": A should remain unchanged")
					}
					assert.Equal(t, tc.expectedFlags.zero, cpu.flags.Zero, tc.name+": Zero flag")
					assert.Equal(t, tc.expectedFlags.carry, cpu.flags.Carry, tc.name+": Carry flag")
					assert.Equal(t, tc.expectedFlags.halfCarry, cpu.flags.HalfCarry, tc.name+": HalfCarry flag")
					assert.Equal(t, tc.expectedFlags.subtract, cpu.flags.Subtract, tc.name+": Subtract flag")
				})
			}
		})
	})

	t.Run("Flow Control", func(t *testing.T) {
		t.Run("Operation: PUSH/POP", func(t *testing.T) {
			t.Run("PUSH_BC", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xC5)
				cpu.registers.b, cpu.registers.c = fromRegisterPair(0x1234)
				cpu.registers.sp = 0xFFFE
				cpu.Step()
				assert.Equal(t, uint16(0xFFFC), cpu.registers.sp, "SP should decrease by 2 after PUSH")
				high := bus.Read(cpu.registers.sp + 1)
				low := bus.Read(cpu.registers.sp)
				assert.Equal(t, uint8(0x12), high, "PUSH_BC: high byte")
				assert.Equal(t, uint8(0x34), low, "PUSH_BC: low byte")
			})

			t.Run("POP_BC", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xC1)
				cpu.registers.sp = 0xFFFC
				bus.Write(cpu.registers.sp, 0x9A)
				bus.Write(cpu.registers.sp+1, 0x78)
				cpu.Step()
				assert.Equal(t, uint16(0xFFFE), cpu.registers.sp, "SP should increase by 2 after POP")
				assert.Equal(t, uint8(0x78), cpu.registers.b, "POP_BC: register B")
				assert.Equal(t, uint8(0x9A), cpu.registers.c, "POP_BC: register C")
			})
		})

		t.Run("Instruction: CALL", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xCD, 0x34, 0x12)
			cpu.registers.sp = 0xFFFE
			initPC := cpu.registers.pc
			cpu.Step()
			assert.Equal(t, uint16(0x1234), cpu.registers.pc, "CALL should jump to target address")
			assert.Equal(t, uint16(0xFFFC), cpu.registers.sp, "CALL should push return address onto stack")
			retLow := bus.Read(cpu.registers.sp)
			retHigh := bus.Read(cpu.registers.sp + 1)
			expectedRet := initPC + 3
			actualRet := toRegisterPair(retHigh, retLow)
			assert.Equal(t, expectedRet, actualRet, "CALL should push correct return address")
		})

		t.Run("Instruction: RET", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xC9)
			cpu.registers.sp = 0xFFFC
			bus.Write(cpu.registers.sp, 0x67)
			bus.Write(cpu.registers.sp+1, 0x45)
			cpu.Step()
			assert.Equal(t, uint16(0x4567), cpu.registers.pc, "RET should set PC to return address")
			assert.Equal(t, uint16(0xFFFE), cpu.registers.sp, "RET should pop return address from stack")
		})

		t.Run("Instruction: SBC_A_d8", func(t *testing.T) {
			type sbcTest struct {
				name         string
				initA        uint8
				immediate    uint8
				initCarry    bool
				expectedA    uint8
				expectedZero bool
			}
			tests := []sbcTest{
				{"SBC_no_carry", 0x05, 0x03, false, 0x02, false},
				{"SBC_with_carry", 0x05, 0x03, true, 0x01, false},
				{"SBC_result_zero", 0x03, 0x02, true, 0x00, true},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, cpu := setupWithOpcode(0xDE, tt.immediate)
					cpu.registers.a = tt.initA
					cpu.flags.Carry = tt.initCarry
					cpu.Step()
					assert.Equal(t, tt.expectedA, cpu.registers.a, "SBC A,d8 result mismatch")
					assert.Equal(t, tt.expectedZero, cpu.flags.Zero, "SBC A,d8 zero flag mismatch")
				})
			}
		})

		t.Run("Instruction: JP_nn", func(t *testing.T) {
			_, cpu := setupWithOpcode(0xC3, 0x21, 0x43)
			cpu.Step()
			assert.Equal(t, uint16(0x4321), cpu.registers.pc, "JP should jump to immediate address")
		})

		t.Run("Instruction: RST", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xDF)
			cpu.registers.sp = 0xFFFE
			initPC := cpu.registers.pc
			cpu.Step()
			assert.Equal(t, uint16(0x0018), cpu.registers.pc, "RST should set PC to fixed vector 0x0018")
			assert.Equal(t, uint16(0xFFFC), cpu.registers.sp, "RST should push return address onto stack")
			retHigh := bus.Read(cpu.registers.sp + 1)
			retLow := bus.Read(cpu.registers.sp)
			expectedRet := initPC + 1
			actualRet := uint16(retHigh)<<8 | uint16(retLow)
			assert.Equal(t, expectedRet, actualRet, "RST should push correct return address")
		})
	})

	t.Run("Conditional Helpers", func(t *testing.T) {
		// Setup a dummy CPU for direct helper calls.
		_, cpu := setupWithOpcode(0x00) // opcode is irrelevant here
		initPC := cpu.registers.pc
		spBefore := cpu.registers.sp

		// Test jump with false condition: should add 3.
		cpu.jump(0x2000, false)
		assert.Equal(t, initPC+3, cpu.registers.pc, "Conditional jump (false) should increment PC by 3")

		// Reset PC.
		cpu.registers.pc = initPC
		// Test jumpRelative with false condition: should add 2.
		var offset int8 = 5
		cpu.jumpRelative(offset, false)
		assert.Equal(t, initPC+2, cpu.registers.pc, "Conditional jumpRelative (false) should increment PC by 2")

		// Test call with false condition: should add 3 and not change SP.
		cpu.registers.pc = initPC
		cpu.call(0x3000, false)
		assert.Equal(t, initPC+3, cpu.registers.pc, "Conditional call (false) should increment PC by 3")
		assert.Equal(t, spBefore, cpu.registers.sp, "Conditional call (false) should not change SP")

		// Test ret with false condition: should add 1.
		cpu.registers.pc = initPC
		cpu.ret(false)
		assert.Equal(t, initPC+1, cpu.registers.pc, "Conditional ret (false) should increment PC by 1")
	})
}
