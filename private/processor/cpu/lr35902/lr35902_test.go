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
	ie := &registers.Interrupt{}
	io := registers.NewRegisters(bus, ie)
	bus.AddDevice(0x0000, 0xFFFF, &memory.Memory{Buffer: make([]byte, 0x10000)})
	cpu := NewLR35902(bus, io, ie)
	// Write provided opcodes to PC sequentially.
	pc := cpu.Registers.PC
	for i, code := range codes {
		bus.Write(pc+uint16(i), code)
	}
	return bus, cpu
}

func TestInstructions(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		t.Run("Instruction: NOP", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x00)
			initPC := cpu.Registers.PC
			cpu.Step()
			assert.Equal(t, initPC+1, cpu.Registers.PC, "NOP should increment PC by 1")
		})

		t.Run("Instruction: LD_BC_d16", func(t *testing.T) {
			// Pass opcode and immediate 16-bit little-endian bytes.
			_, cpu := setupWithOpcode(0x01, 0x42, 0x24)
			cpu.Step()
			assert.Equal(t, uint16(0x2442), toRegisterPair(cpu.Registers.B, cpu.Registers.C), "BC should load immediate 16-bit value")
		})
	})

	t.Run("8-Bit Loads", func(t *testing.T) {
		t.Run("Instruction: LD_d8", func(t *testing.T) {
			{
				_, cpu := setupWithOpcode(0x06, 0x42)
				cpu.Step()
				assert.Equal(t, uint8(0x42), cpu.Registers.B, "B should load immediate 8-bit value")
			}
			{
				_, cpu := setupWithOpcode(0x0E, 0x55)
				cpu.Step()
				assert.Equal(t, uint8(0x55), cpu.Registers.C, "C should load immediate 8-bit value")
			}
		})
		// Example for an operation test:
		t.Run("Operation: INC B", func(t *testing.T) {
			{
				_, cpu := setupWithOpcode(0x04)
				cpu.Registers.B = 1
				cpu.Step()
				assert.Equal(t, uint8(2), cpu.Registers.B, "B should increment by 1")
				assert.False(t, cpu.Flags.Zero, "Z flag should be reset on INC")
				assert.False(t, cpu.Flags.Subtract, "N flag should be reset on INC")
				assert.False(t, cpu.Flags.HalfCarry, "H flag should be reset on INC")
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
					cpu.Registers.B = tt.initVal
					cpu.Step()
					assert.Equal(t, tt.expResult, cpu.Registers.B, "DEC B did not produce expected value")
					assert.Equal(t, tt.expZero, cpu.Flags.Zero, "Zero flag mismatch on DEC B")
					assert.True(t, cpu.Flags.Subtract, "N flag should be set on DEC")
				})
			}
		})
	})

	t.Run("16-Bit Operations", func(t *testing.T) {
		t.Run("Instruction: INC/DEC BC", func(t *testing.T) {
			{
				_, cpu := setupWithOpcode(0x03)
				cpu.Registers.B, cpu.Registers.C = fromRegisterPair(0x01)
				cpu.Step()
				assert.Equal(t, uint16(2), toRegisterPair(cpu.Registers.B, cpu.Registers.C), "BC should increment by 1")
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
					cpu.Registers.B, cpu.Registers.C = fromRegisterPair(tt.initBC)
					cpu.Step()
					res := toRegisterPair(cpu.Registers.B, cpu.Registers.C)
					assert.Equal(t, tt.expResult, res, "DEC BC did not produce expected value")
				})
			}
		})
		t.Run("Operation: ADD HL,BC", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x09)
			cpu.Registers.H, cpu.Registers.L = fromRegisterPair(0x1)
			cpu.Registers.B, cpu.Registers.C = fromRegisterPair(0x1)
			cpu.Step()
			assert.Equal(t, uint16(0x2), toRegisterPair(cpu.Registers.H, cpu.Registers.L), "HL should add BC")
			assert.False(t, cpu.Flags.Subtract, "N flag should be reset in addition")
		})
		t.Run("ADD HL,BC bits ordering", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x09)
			// Set HL = 0x1234 and BC = 0x4321.
			cpu.Registers.H, cpu.Registers.L = fromRegisterPair(0x1234)
			cpu.Registers.B, cpu.Registers.C = fromRegisterPair(0x4321)
			cpu.Step()
			expectedHL := uint16(0x1234 + 0x4321) // should equal 0x5555
			actualHL := toRegisterPair(cpu.Registers.H, cpu.Registers.L)
			if actualHL != expectedHL {
				t.Errorf("ADD HL,BC produced 0x%04X; expected 0x%04X", actualHL, expectedHL)
			}
		})
	})

	t.Run("Memory", func(t *testing.T) {
		t.Run("Instruction: Address Operations", func(t *testing.T) {
			t.Run("LD_A_from_BC", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0x02)
				cpu.Registers.B, cpu.Registers.C = fromRegisterPair(0x01)
				cpu.Registers.A = 0xAA
				cpu.Step()
				assert.Equal(t, uint8(0xAA), bus.Read(0x0001), "Memory at address BC should be loaded with A")
			})
			t.Run("LD_A_from_a16", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0x0A)
				cpu.Registers.B = 0x01
				cpu.Registers.C = 0x00
				bus.Write(0x100, 0xBB)
				cpu.Step()
				assert.Equal(t, uint8(0xBB), cpu.Registers.A, "A should load value from memory at address BC")
			})
			t.Run("LD_a16_from_A", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0x08, 0x0B, 0x00)
				cpu.Registers.SP = 0x1234
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
				bus.Write(cpu.Registers.PC+1, uint8(addr&0xFF))
				bus.Write(cpu.Registers.PC+2, uint8(addr>>8))
				cpu.Step()
				assert.Equal(t, uint8(0x7F), cpu.Registers.A, "LD A,(a16) should load value from memory")
			})
			t.Run("LD_a16_from_A", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xEA)
				addr := uint16(0x3000)
				cpu.Registers.A = 0x3C
				bus.Write(cpu.Registers.PC+1, uint8(addr&0xFF))
				bus.Write(cpu.Registers.PC+2, uint8(addr>>8))
				cpu.Step()
				assert.Equal(t, uint8(0x3C), bus.Read(addr), "LD (a16),A should store A into memory")
			})
			t.Run("LDH_A_from_n", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xF0)
				offset := uint8(0x20)
				addr := uint16(0xFF00) + uint16(offset)
				bus.Write(addr, 0x99)
				bus.Write(cpu.Registers.PC+1, offset)
				cpu.Step()
				assert.Equal(t, uint8(0x99), cpu.Registers.A, "LDH A,(n) should load value from 0xFF00+n")
			})
			t.Run("LDH_n_from_A", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xE0)
				offset := uint8(0x30)
				cpu.Registers.A = 0xAB
				addr := uint16(0xFF00) + uint16(offset)
				bus.Write(cpu.Registers.PC+1, offset)
				cpu.Step()
				assert.Equal(t, uint8(0xAB), bus.Read(addr), "LDH (n),A should store A into memory at 0xFF00+n")
			})
		})
	})

	t.Run("Rotation", func(t *testing.T) {
		t.Run("Instruction: RLCA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x07)
			cpu.Registers.A = 0x80
			cpu.Step()
			assert.Equal(t, uint8(0x01), cpu.Registers.A, "A should rotate left (RLCA)")
			assert.True(t, cpu.Flags.Carry, "Carry flag should be set by RLCA")
		})
		t.Run("Instruction: RRCA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x0F)
			cpu.Registers.A = 0x01
			cpu.Step()
			assert.Equal(t, uint8(0x80), cpu.Registers.A, "A should rotate right (RRCA)")
			assert.True(t, cpu.Flags.Carry, "Carry flag should be set by RRCA")
		})
		t.Run("Instruction: RLA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x17)
			cpu.Registers.A = 0x80
			cpu.Flags.Carry = true
			cpu.Step()
			assert.Equal(t, uint8(0x01), cpu.Registers.A, "A should rotate left (RLA)")
			assert.True(t, cpu.Flags.Carry, "Carry flag should be set by RLA")
		})
		t.Run("Instruction: RRA", func(t *testing.T) {
			_, cpu := setupWithOpcode(0x1F)
			cpu.Registers.A = 0x01
			cpu.Flags.Carry = true
			cpu.Step()
			assert.Equal(t, uint8(0x80), cpu.Registers.A, "A should rotate right (RRA)")
			// Updated: Carry flag should remain true after RRA.
			assert.True(t, cpu.Flags.Carry, "Carry flag should be set by RRA")
		})

		t.Run("CB Rotation", func(t *testing.T) {
			t.Run("Instruction: RLC_B", func(t *testing.T) {
				_, cpu := setupWithOpcode(0xCB, 0x00)
				cpu.Registers.B = 0x80 // 10000000
				cpu.Step()
				cpu.Step()
				assert.Equal(t, uint8(0x01), cpu.Registers.B, "RLC_B should rotate B left, result 0x01")
				assert.True(t, cpu.Flags.Carry, "RLC_B should set Carry flag")
			})

			t.Run("Instruction: RRC_B", func(t *testing.T) {
				_, cpu := setupWithOpcode(0xCB, 0x08)
				cpu.Registers.B = 0x01 // 00000001
				cpu.Step()
				cpu.Step()
				assert.Equal(t, uint8(0x80), cpu.Registers.B, "RRC_B should rotate B right, result 0x80")
				assert.True(t, cpu.Flags.Carry, "RRC_B should set Carry flag")
			})

			t.Run("Instruction: RL_B", func(t *testing.T) {
				_, cpu := setupWithOpcode(0xCB, 0x10)
				cpu.Registers.B = 0x80 // 10000000
				cpu.Flags.Carry = true // initial carry is set
				cpu.Step()
				cpu.Step()
				assert.Equal(t, uint8(0x01), cpu.Registers.B, "RL_B should rotate B left through carry, result 0x01")
				assert.True(t, cpu.Flags.Carry, "RL_B should set Carry flag")
			})

			t.Run("Instruction: RR_B", func(t *testing.T) {
				_, cpu := setupWithOpcode(0xCB, 0x18)
				cpu.Registers.B = 0x01 // 00000001
				cpu.Flags.Carry = true // initial carry is set
				cpu.Step()
				cpu.Step()
				assert.Equal(t, uint8(0x80), cpu.Registers.B, "RR_B should rotate B right through carry, result 0x80")
				assert.True(t, cpu.Flags.Carry, "RR_B should set Carry flag")
			})

			t.Run("Instruction: RLC_(HL)", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xCB, 0x06)
				hlAddr := uint16(0x2000)
				cpu.Registers.H, cpu.Registers.L = fromRegisterPair(hlAddr)
				// Write initial value: 0x85 (10000101)
				bus.Write(hlAddr, 0x85)
				cpu.Step()
				cpu.Step()
				// Expected: (0x85<<1 | (0x85>>7)) = (0x0A | 0x01) = 0x0B
				assert.Equal(t, uint8(0x0B), bus.Read(hlAddr), "RLC_(HL) should rotate value in memory at HL")
				assert.True(t, cpu.Flags.Carry, "RLC_(HL) should set Carry flag")
			})

			t.Run("Instruction: RRC_(HL)", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xCB, 0x0E)
				hlAddr := uint16(0x2000)
				cpu.Registers.H, cpu.Registers.L = fromRegisterPair(hlAddr)
				// Write initial value: 0x01 (00000001)
				bus.Write(hlAddr, 0x01)
				cpu.Step()
				cpu.Step()
				// Expected: (0x01>>1 | (0x01<<7)&0xFF) = (0x00 | 0x80) = 0x80
				assert.Equal(t, uint8(0x80), bus.Read(hlAddr), "RRC_(HL) should rotate value in memory at HL")
				assert.True(t, cpu.Flags.Carry, "RRC_(HL) should set Carry flag")
			})

			t.Run("Instruction: SRL A", func(t *testing.T) {
				_, cpu := setupWithOpcode(0xCB, 0x3F) // CB prefix, SRL A opcode
				cpu.Registers.A = 0x02                // binary: 0000 0010
				cpu.Step()                            // process CB prefix
				cpu.Step()                            // execute SRL A
				assert.Equal(t, uint8(0x01), cpu.Registers.A, "SRL A: A should be shifted right logically")
				assert.False(t, cpu.Flags.Zero, "SRL A: Zero flag should be false")
				assert.False(t, cpu.Flags.HalfCarry, "SRL A: HalfCarry flag should be false")
				assert.False(t, cpu.Flags.Subtract, "SRL A: Subtract flag should be false")
				assert.False(t, cpu.Flags.Carry, "SRL A: Carry flag should reflect LSB (expected false)")
			})

			t.Run("Instruction: SRL (HL)", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xCB, 0x3E) // CB prefix, SRL (HL) opcode
				hlAddr := uint16(0x2000)
				cpu.Registers.H, cpu.Registers.L = fromRegisterPair(hlAddr)
				bus.Write(hlAddr, 0x03) // binary: 0000 0011; expected result: 0x01 with Carry true
				cpu.Step()              // process CB prefix
				cpu.Step()              // execute SRL (HL)
				assert.Equal(t, uint8(0x01), bus.Read(hlAddr), "SRL (HL): value should be shifted right logically")
				assert.False(t, cpu.Flags.Zero, "SRL (HL): Zero flag should be false")
				assert.False(t, cpu.Flags.HalfCarry, "SRL (HL): HalfCarry flag should be false")
				assert.False(t, cpu.Flags.Subtract, "SRL (HL): Subtract flag should be false")
				assert.True(t, cpu.Flags.Carry, "SRL (HL): Carry flag should be set (LSB was 1)")
			})
		})
	})

	t.Run("CB Bit, Res and Set", func(t *testing.T) {
		// BIT tests on register B
		t.Run("BIT 0, B - bit set", func(t *testing.T) {
			_, cpu := setupWithOpcode(0xCB, 0x40) // BIT 0, B
			cpu.Registers.B = 0x01                // bit0 is set
			cpu.Flags.Carry = true                // initial carry value
			cpu.Step()                            // process CB prefix
			cpu.Step()                            // execute BIT 0, B
			// Expected: bit is set -> Zero false, H set, N reset, Carry unchanged.
			assert.Equal(t, false, cpu.Flags.Zero, "BIT 0, B: Zero flag should be reset when bit is set")
			assert.Equal(t, true, cpu.Flags.HalfCarry, "BIT 0, B: HalfCarry flag should be set")
			assert.Equal(t, false, cpu.Flags.Subtract, "BIT 0, B: Subtract flag should be reset")
			assert.Equal(t, true, cpu.Flags.Carry, "BIT 0, B: Carry flag should remain unchanged")
		})

		t.Run("BIT 0, B - bit clear", func(t *testing.T) {
			_, cpu := setupWithOpcode(0xCB, 0x40) // BIT 0, B
			cpu.Registers.B = 0x00                // bit0 clear
			cpu.Flags.Carry = false
			cpu.Step()
			cpu.Step()
			// Expected: bit is clear -> Zero true, H set, N reset, Carry unchanged.
			assert.Equal(t, true, cpu.Flags.Zero, "BIT 0, B: Zero flag should be set when bit is clear")
			assert.Equal(t, true, cpu.Flags.HalfCarry, "BIT 0, B: HalfCarry flag should be set")
			assert.Equal(t, false, cpu.Flags.Subtract, "BIT 0, B: Subtract flag should be reset")
			assert.Equal(t, false, cpu.Flags.Carry, "BIT 0, B: Carry flag should remain unchanged")
		})

		// BIT test on (HL)
		t.Run("BIT 3, (HL) - bit set", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xCB, 0x5E) // BIT 3, (HL)
			hlAddr := uint16(0x2000)
			cpu.Registers.H, cpu.Registers.L = fromRegisterPair(hlAddr)
			bus.Write(hlAddr, 0x08) // 0x08 has bit3 set (0000 1000)
			cpu.Flags.Carry = false
			cpu.Step()
			cpu.Step()
			// Expected: bit is set -> Zero false, H set, N reset.
			assert.Equal(t, false, cpu.Flags.Zero, "BIT 3,(HL): Zero flag should be reset when bit is set")
			assert.Equal(t, true, cpu.Flags.HalfCarry, "BIT 3,(HL): HalfCarry flag should be set")
			assert.Equal(t, false, cpu.Flags.Subtract, "BIT 3,(HL): Subtract flag should be reset")
		})

		// RES tests
		t.Run("RES 0, B", func(t *testing.T) {
			_, cpu := setupWithOpcode(0xCB, 0x80) // RES 0, B
			cpu.Registers.B = 0xFF                // all bits set
			cpu.Step()
			cpu.Step()
			// Expected: reset bit0 -> 0xFE.
			assert.Equal(t, uint8(0xFE), cpu.Registers.B, "RES 0, B should reset bit 0")
		})

		t.Run("RES 1, (HL)", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xCB, 0x8E) // RES 1, (HL)
			hlAddr := uint16(0x2000)
			cpu.Registers.H, cpu.Registers.L = fromRegisterPair(hlAddr)
			bus.Write(hlAddr, 0xFF) // all bits set
			cpu.Step()
			cpu.Step()
			// Expected: reset bit1 -> 0xFD (1111 1101).
			assert.Equal(t, uint8(0xFD), bus.Read(hlAddr), "RES 1,(HL) should reset bit 1")
		})

		// SET tests
		t.Run("SET 0, B", func(t *testing.T) {
			_, cpu := setupWithOpcode(0xCB, 0xC0) // SET 0, B
			cpu.Registers.B = 0xFE                // bit0 is clear
			cpu.Step()
			cpu.Step()
			// Expected: set bit0 -> 0xFF.
			assert.Equal(t, uint8(0xFF), cpu.Registers.B, "SET 0, B should set bit 0")
		})

		t.Run("SET 1, (HL)", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xCB, 0xCE) // SET 1, (HL)
			hlAddr := uint16(0x2000)
			cpu.Registers.H, cpu.Registers.L = fromRegisterPair(hlAddr)
			bus.Write(hlAddr, 0xF9) // 0xF9: 1111 1001, bit1 clear
			cpu.Step()
			cpu.Step()
			// Expected: set bit1 -> 0xFB (1111 1011).
			assert.Equal(t, uint8(0xFB), bus.Read(hlAddr), "SET 1,(HL) should set bit 1")
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
					cpu.Registers.A = tc.initA
					cpu.Registers.B = tc.initB
					cpu.Flags.Carry = tc.initCarry
					cpu.Step()
					if !tc.checkAUnchanged {
						assert.Equal(t, tc.expectedA, cpu.Registers.A, tc.name+": A value")
					} else {
						assert.Equal(t, tc.initA, cpu.Registers.A, tc.name+": A should remain unchanged")
					}
					assert.Equal(t, tc.expectedFlags.zero, cpu.Flags.Zero, tc.name+": Zero flag")
					assert.Equal(t, tc.expectedFlags.carry, cpu.Flags.Carry, tc.name+": Carry flag")
					assert.Equal(t, tc.expectedFlags.halfCarry, cpu.Flags.HalfCarry, tc.name+": HalfCarry flag")
					assert.Equal(t, tc.expectedFlags.subtract, cpu.Flags.Subtract, tc.name+": Subtract flag")
				})
			}
		})
		t.Run("Instruction: SBC A,B", func(t *testing.T) {
			// Table-driven tests for SBC A,B (opcode 0x98)
			type sbcTest struct {
				name         string
				initA        uint8
				initB        uint8
				initCarry    bool
				expectedA    uint8
				expZero      bool
				expHalfCarry bool
				expCarry     bool
			}
			tests := []sbcTest{
				{"SBC_A_B_no_borrow", 0x05, 0x03, false, 0x02, false, false, false},
				{"SBC_A_B_with_borrow", 0x05, 0x03, true, 0x01, false, false, false},
				{"SBC_A_B_result_zero", 0x03, 0x03, false, 0x00, true, false, false},
				{"SBC_A_B_underflow", 0x00, 0x01, false, 0xFF, false, true, true},
			}
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					_, cpu := setupWithOpcode(0x98) // SBC A,B opcode
					cpu.Registers.A = tc.initA
					cpu.Registers.B = tc.initB
					cpu.Flags.Carry = tc.initCarry
					cpu.Step()
					assert.Equal(t, tc.expectedA, cpu.Registers.A, tc.name+": A value mismatch")
					assert.Equal(t, tc.expZero, cpu.Flags.Zero, tc.name+": Zero flag mismatch")
					// For subtraction, N flag is always set.
					assert.True(t, cpu.Flags.Subtract, tc.name+": Subtract flag should be set")
					assert.Equal(t, tc.expHalfCarry, cpu.Flags.HalfCarry, tc.name+": HalfCarry flag mismatch")
					assert.Equal(t, tc.expCarry, cpu.Flags.Carry, tc.name+": Carry flag mismatch")
				})
			}
		})
	})

	t.Run("Flow Control", func(t *testing.T) {
		t.Run("Operation: PUSH/POP", func(t *testing.T) {
			t.Run("PUSH_BC", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xC5)
				cpu.Registers.B, cpu.Registers.C = fromRegisterPair(0x1234)
				cpu.Registers.SP = 0xFFFE
				cpu.Step()
				assert.Equal(t, uint16(0xFFFC), cpu.Registers.SP, "SP should decrease by 2 after PUSH")
				high := bus.Read(cpu.Registers.SP + 1)
				low := bus.Read(cpu.Registers.SP)
				assert.Equal(t, uint8(0x12), high, "PUSH_BC: high byte")
				assert.Equal(t, uint8(0x34), low, "PUSH_BC: low byte")
			})

			t.Run("POP_BC", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xC1)
				cpu.Registers.SP = 0xFFFC
				bus.Write(cpu.Registers.SP, 0x9A)
				bus.Write(cpu.Registers.SP+1, 0x78)
				cpu.Step()
				assert.Equal(t, uint16(0xFFFE), cpu.Registers.SP, "SP should increase by 2 after POP")
				assert.Equal(t, uint8(0x78), cpu.Registers.B, "POP_BC: register B")
				assert.Equal(t, uint8(0x9A), cpu.Registers.C, "POP_BC: register C")
			})
		})

		t.Run("Instruction: CALL", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xCD, 0x34, 0x12)
			cpu.Registers.SP = 0xFFFE
			initPC := cpu.Registers.PC
			cpu.Step()
			assert.Equal(t, uint16(0x1234), cpu.Registers.PC, "CALL should jump to target address")
			assert.Equal(t, uint16(0xFFFC), cpu.Registers.SP, "CALL should push return address onto stack")
			retLow := bus.Read(cpu.Registers.SP)
			retHigh := bus.Read(cpu.Registers.SP + 1)
			expectedRet := initPC + 3
			actualRet := toRegisterPair(retHigh, retLow)
			assert.Equal(t, expectedRet, actualRet, "CALL should push correct return address")
		})

		t.Run("Instruction: RET", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xC9)
			cpu.Registers.SP = 0xFFFC
			bus.Write(cpu.Registers.SP, 0x67)
			bus.Write(cpu.Registers.SP+1, 0x45)
			cpu.Step()
			assert.Equal(t, uint16(0x4567), cpu.Registers.PC, "RET should set PC to return address")
			assert.Equal(t, uint16(0xFFFE), cpu.Registers.SP, "RET should pop return address from stack")
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
					cpu.Registers.A = tt.initA
					cpu.Flags.Carry = tt.initCarry
					cpu.Step()
					assert.Equal(t, tt.expectedA, cpu.Registers.A, "SBC A,d8 result mismatch")
					assert.Equal(t, tt.expectedZero, cpu.Flags.Zero, "SBC A,d8 zero flag mismatch")
				})
			}
		})

		t.Run("Instruction: JP_nn", func(t *testing.T) {
			_, cpu := setupWithOpcode(0xC3, 0x21, 0x43)
			cpu.Step()
			assert.Equal(t, uint16(0x4321), cpu.Registers.PC, "JP should jump to immediate address")
		})

		t.Run("Instruction: RST", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xDF)
			cpu.Registers.SP = 0xFFFE
			initPC := cpu.Registers.PC
			cpu.Step()
			assert.Equal(t, uint16(0x0018), cpu.Registers.PC, "RST should set PC to fixed vector 0x0018")
			assert.Equal(t, uint16(0xFFFC), cpu.Registers.SP, "RST should push return address onto stack")
			retHigh := bus.Read(cpu.Registers.SP + 1)
			retLow := bus.Read(cpu.Registers.SP)
			expectedRet := initPC + 1
			actualRet := uint16(retHigh)<<8 | uint16(retLow)
			assert.Equal(t, expectedRet, actualRet, "RST should push correct return address")
		})

		// New test for LD HL,SP+r8 (opcode 0xF8)
		t.Run("Instruction: LD HL,SP+r8", func(t *testing.T) {
			tests := []struct {
				name       string
				sp         uint16
				offset     int8
				expectedHL uint16
				expH       bool
				expC       bool
			}{
				{"no flags", 0x1000, 5, 0x1005, false, false},
				{"half carry", 0x0015, 0x0B, 0x0020, true, false},
				{"carry", 0x00F0, 0x14, 0x0104, false, true},
				// Updated: For negative offset, computed via unsigned arithmetic:
				{"negative offset", 0x1000, -3, 0x0FFD, false, false},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, cpu := setupWithOpcode(0xF8, uint8(tt.offset))
					cpu.Registers.SP = tt.sp
					cpu.Step()
					hl := toRegisterPair(cpu.Registers.H, cpu.Registers.L)
					assert.Equal(t, tt.expectedHL, hl, tt.name+": HL mismatch")
					assert.False(t, cpu.Flags.Zero, tt.name+": Zero flag must be reset")
					assert.False(t, cpu.Flags.Subtract, tt.name+": Subtract flag must be reset")
					assert.Equal(t, tt.expH, cpu.Flags.HalfCarry, tt.name+": HalfCarry flag mismatch")
					assert.Equal(t, tt.expC, cpu.Flags.Carry, tt.name+": Carry flag mismatch")
				})
			}
		})

		t.Run("POP_AF", func(t *testing.T) {
			t.Run("all flags set", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xF1)
				cpu.Registers.SP = 0xFFFC
				// Write: low byte (flags) = 0xF0 (11110000), high byte (A) = 0xAA
				bus.Write(cpu.Registers.SP, 0xF0)
				bus.Write(cpu.Registers.SP+1, 0xAA)
				cpu.popAF()
				assert.Equal(t, uint8(0xAA), cpu.Registers.A, "POP_AF: A register mismatch")
				assert.True(t, cpu.Flags.Zero, "POP_AF: Zero flag should be set")
				assert.True(t, cpu.Flags.Subtract, "POP_AF: Subtract flag should be set")
				assert.True(t, cpu.Flags.HalfCarry, "POP_AF: HalfCarry flag should be set")
				assert.True(t, cpu.Flags.Carry, "POP_AF: Carry flag should be set")
			})
			t.Run("no flags set", func(t *testing.T) {
				bus, cpu := setupWithOpcode(0xF1)
				cpu.Registers.SP = 0xFFFC
				// Write: low byte (flags) = 0x00, high byte (A) = 0x55
				bus.Write(cpu.Registers.SP, 0x00)
				bus.Write(cpu.Registers.SP+1, 0x55)
				cpu.popAF()
				assert.Equal(t, uint8(0x55), cpu.Registers.A, "POP_AF: A register mismatch")
				assert.False(t, cpu.Flags.Zero, "POP_AF: Zero flag should be reset")
				assert.False(t, cpu.Flags.Subtract, "POP_AF: Subtract flag should be reset")
				assert.False(t, cpu.Flags.HalfCarry, "POP_AF: HalfCarry flag should be reset")
				assert.False(t, cpu.Flags.Carry, "POP_AF: Carry flag should be reset")
			})
		})

		t.Run("Instruction: ADD SP,r8", func(t *testing.T) {
			tests := []struct {
				name       string
				sp         uint16
				offset     int8
				expectedSP uint16
				expH       bool
				expC       bool
			}{
				{"no flags", 0x1000, 5, 0x1005, false, false},
				{"half carry", 0x100F, 1, 0x1010, true, false},
				{"carry", 0x10FF, 1, 0x1100, true, true},
				{"negative offset", 0x1000, -3, 0x0FFD, false, false},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, cpu := setupWithOpcode(0xE8, uint8(tt.offset))
					cpu.Registers.SP = tt.sp
					cpu.Step()
					assert.Equal(t, tt.expectedSP, cpu.Registers.SP, tt.name+": SP mismatch")
					// Flags: Z and N are reset; check HalfCarry and Carry.
					assert.False(t, cpu.Flags.Zero, tt.name+": Zero flag must be reset")
					assert.False(t, cpu.Flags.Subtract, tt.name+": Subtract flag must be reset")
					assert.Equal(t, tt.expH, cpu.Flags.HalfCarry, tt.name+": HalfCarry flag mismatch")
					assert.Equal(t, tt.expC, cpu.Flags.Carry, tt.name+": Carry flag mismatch")
				})
			}
		})

		t.Run("Instruction: EI delay", func(t *testing.T) {
			// Set up CPU with EI (0xFB) followed by a NOP (0x00)
			_, cpu := setupWithOpcode(0xFB, 0x00)
			// Ensure initial IME is false.
			cpu.ime = false

			// Execute EI; IME should remain false immediately due to the EI delay.
			cpu.Step()
			if cpu.ime {
				t.Error("After EI instruction, IME should not be enabled immediately")
			}

			// Execute the following NOP instruction; now IME should be enabled.
			cpu.Step()
			if !cpu.ime {
				t.Error("After one instruction delay, IME should be enabled")
			}
		})

		t.Run("Instruction: RETI", func(t *testing.T) {
			bus, cpu := setupWithOpcode(0xD9)
			cpu.Registers.SP = 0xFFFC
			bus.Write(cpu.Registers.SP, 0x34)   // low byte
			bus.Write(cpu.Registers.SP+1, 0x12) // high byte => target address 0x1234
			cpu.Step()
			assert.Equal(t, uint16(0x1234), cpu.Registers.PC, "RETI should set PC from return address")
			assert.Equal(t, uint16(0xFFFE), cpu.Registers.SP, "RETI should pop return address from stack")
			assert.True(t, cpu.ime, "RETI should set IME to true")
		})
	})

	t.Run("Conditional Helpers", func(t *testing.T) {
		// Setup a dummy CPU for direct helper calls.
		_, cpu := setupWithOpcode(0x00) // opcode is irrelevant here
		initPC := cpu.Registers.PC
		spBefore := cpu.Registers.SP

		// Test jump with false condition: should add 3.
		cpu.jump(0x2000, false)
		assert.Equal(t, initPC+3, cpu.Registers.PC, "Conditional jump (false) should increment PC by 3")

		// Reset PC.
		cpu.Registers.PC = initPC
		// Test jumpRelative with false condition: should add 2.
		var offset int8 = 5
		cpu.jumpRelative(offset, false)
		assert.Equal(t, initPC+2, cpu.Registers.PC, "Conditional jumpRelative (false) should increment PC by 2")

		// Test call with false condition: should add 3 and not change SP.
		cpu.Registers.PC = initPC
		cpu.call(0x3000, false)
		assert.Equal(t, initPC+3, cpu.Registers.PC, "Conditional call (false) should increment PC by 3")
		assert.Equal(t, spBefore, cpu.Registers.SP, "Conditional call (false) should not change SP")

		// Test ret with false condition: should add 1.
		cpu.Registers.PC = initPC
		cpu.ret(false)
		assert.Equal(t, initPC+1, cpu.Registers.PC, "Conditional ret (false) should increment PC by 1")
	})
}
