package lr35902

import (
	"fmt"
	"testing"

	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/stretchr/testify/assert"
)

// setupCPU initializes the memory bus and CPU for each test
func setupCPU() (*memory.Bus, *LR35902) {
	bus := &memory.Bus{}
	io := &registers.Registers{}
	bus.AddDevice(0x0000, 0xFFFF, &memory.Memory{Buffer: make([]byte, 0x10000)})
	cpu := NewLR35902(bus, io)
	return bus, cpu
}

func TestInstructions(t *testing.T) {
	// Group: Basic Instructions
	t.Run("Basic", func(t *testing.T) {
		t.Run("NOP", func(t *testing.T) {
			bus, cpu := setupCPU()
			bus.Write(cpu.registers.pc, 0x00) // NOP opcode
			initPC := cpu.registers.pc
			cpu.Clock()
			assert.Equal(t, initPC+1, cpu.registers.pc, "NOP should increment PC by 1")
		})

		t.Run("LD_BC_d16", func(t *testing.T) {
			bus, cpu := setupCPU()
			bus.Write(cpu.registers.pc, 0x01)   // LD BC,d16 opcode
			bus.Write(cpu.registers.pc+1, 0x42) // Low byte
			bus.Write(cpu.registers.pc+2, 0x24) // High byte
			cpu.Clock()
			assert.Equal(t, uint16(0x2442), toRegisterPair(cpu.registers.b, cpu.registers.c), "BC should load immediate 16-bit value")
		})
	})

	// Group: 8-bit Loads and Increment/Decrement
	t.Run("8bit", func(t *testing.T) {
		t.Run("LD_d8_and_INC_DEC_8bit", func(t *testing.T) {
			// LD B,d8 (opcode 0x06)
			{
				bus, cpu := setupCPU()
				bus.Write(cpu.registers.pc, 0x06)
				bus.Write(cpu.registers.pc+1, 0x42)
				cpu.Clock()
				assert.Equal(t, uint8(0x42), cpu.registers.b, "B should load immediate 8-bit value")
			}
			// LD C,d8 (opcode 0x0E)
			{
				bus, cpu := setupCPU()
				bus.Write(cpu.registers.pc, 0x0E)
				bus.Write(cpu.registers.pc+1, 0x55)
				cpu.Clock()
				assert.Equal(t, uint8(0x55), cpu.registers.c, "C should load immediate 8-bit value")
			}
			// INC B (opcode 0x04)
			{
				bus, cpu := setupCPU()
				cpu.registers.b = 1
				bus.Write(cpu.registers.pc, 0x04)
				cpu.Clock()
				assert.Equal(t, uint8(2), cpu.registers.b, "B should increment by 1")
				assert.False(t, cpu.flags.Zero, "Z flag should be reset on INC")
				assert.False(t, cpu.flags.Subtract, "N flag should be reset on INC")
				assert.False(t, cpu.flags.HalfCarry, "H flag should be reset on INC")
			}

			// Table-driven tests for DEC B
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
					bus, cpu := setupCPU()
					cpu.registers.b = tt.initVal
					bus.Write(cpu.registers.pc, 0x05) // DEC B opcode
					cpu.Clock()
					assert.Equal(t, tt.expResult, cpu.registers.b, "DEC B did not produce expected value")
					assert.Equal(t, tt.expZero, cpu.flags.Zero, "Zero flag mismatch on DEC B")
					assert.True(t, cpu.flags.Subtract, "N flag should be set on DEC")
				})
			}
		})
	})

	// Group: 16-bit Operations
	t.Run("16bit", func(t *testing.T) {
		t.Run("INC_DEC_16bit", func(t *testing.T) {
			// INC BC (opcode 0x03)
			{
				bus, cpu := setupCPU()
				cpu.registers.b = 0x01
				cpu.registers.c = 0x00
				bus.Write(cpu.registers.pc, 0x03)
				cpu.Clock()
				assert.Equal(t, uint16(2), toRegisterPair(cpu.registers.b, cpu.registers.c), "BC should increment by 1")
			}

			// Table-driven tests for DEC BC
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
					bus, cpu := setupCPU()
					// Split 16-bit value into registers b and c.
					cpu.registers.b = uint8(tt.initBC >> 8)
					cpu.registers.c = uint8(tt.initBC & 0xFF)
					bus.Write(cpu.registers.pc, 0x0B) // DEC BC opcode
					cpu.Clock()
					res := toRegisterPair(cpu.registers.b, cpu.registers.c)
					assert.Equal(t, tt.expResult, res, "DEC BC did not produce expected value")
				})
			}
		})

		t.Run("ADD_HL_BC", func(t *testing.T) {
			bus, cpu := setupCPU()
			cpu.registers.h = 0x01
			cpu.registers.l = 0x00
			cpu.registers.b = 0x01
			cpu.registers.c = 0x00
			bus.Write(cpu.registers.pc, 0x09)
			cpu.Clock()
			assert.Equal(t, uint16(0x2), toRegisterPair(cpu.registers.h, cpu.registers.l), "HL should add BC")
			assert.False(t, cpu.flags.Subtract, "N flag should be reset in addition")
		})
	})

	// Group: Memory and Address Operations
	t.Run("Memory", func(t *testing.T) {
		t.Run("Address_Operations", func(t *testing.T) {
			// LD (BC),A (opcode 0x02)
			{
				bus, cpu := setupCPU()
				cpu.registers.b = 0x01
				cpu.registers.c = 0x00
				cpu.registers.a = 0xAA
				bus.Write(cpu.registers.pc, 0x02)
				cpu.Clock()
				assert.Equal(t, uint8(0xAA), bus.Read(0x0100), "Memory at address BC should be loaded with A")
			}
			// LD A,(BC) (opcode 0x0A)
			{
				bus, cpu := setupCPU()
				cpu.registers.b = 0x01
				cpu.registers.c = 0x00
				bus.Write(0x100, 0xBB)
				bus.Write(cpu.registers.pc, 0x0A)
				cpu.Clock()
				assert.Equal(t, uint8(0xBB), cpu.registers.a, "A should load value from memory at address BC")
			}
			// LD a16,SP (opcode 0x08)
			{
				bus, cpu := setupCPU()
				cpu.registers.sp = 0x1234
				bus.Write(cpu.registers.pc, 0x08)
				bus.Write(cpu.registers.pc+1, 0x0B) // target address low
				bus.Write(cpu.registers.pc+2, 0x00) // target address high
				cpu.Clock()
				assert.Equal(t, uint8(0x34), bus.Read(0x0B), "Memory low should be SP's low byte")
				assert.Equal(t, uint8(0x12), bus.Read(0x0C), "Memory high should be SP's high byte")
			}
		})

		t.Run("LD_variants", func(t *testing.T) {
			// Test LD A,(a16): opcode 0xFA
			t.Run("LD_A_from_a16", func(t *testing.T) {
				bus, cpu := setupCPU()
				// Set target memory value at address 0x2000.
				addr := uint16(0x2000)
				bus.Write(addr, 0x7F)
				// Write opcode and a16 address (little-endian).
				bus.Write(cpu.registers.pc, 0xFA)
				bus.Write(cpu.registers.pc+1, uint8(addr&0xFF))
				bus.Write(cpu.registers.pc+2, uint8(addr>>8))
				cpu.Clock()
				assert.Equal(t, uint8(0x7F), cpu.registers.a, "LD A,(a16) should load value from memory")
			})

			// Test LD (a16),A: opcode 0xEA
			t.Run("LD_a16_from_A", func(t *testing.T) {
				bus, cpu := setupCPU()
				cpu.registers.a = 0x3C
				addr := uint16(0x3000)
				bus.Write(cpu.registers.pc, 0xEA)
				bus.Write(cpu.registers.pc+1, uint8(addr&0xFF))
				bus.Write(cpu.registers.pc+2, uint8(addr>>8))
				cpu.Clock()
				assert.Equal(t, uint8(0x3C), bus.Read(addr), "LD (a16),A should store A into memory")
			})

			// Test LDH A,(n): opcode 0xF0, loads from 0xFF00+n.
			t.Run("LDH_A_from_n", func(t *testing.T) {
				bus, cpu := setupCPU()
				// Use immediate offset 0x20.
				offset := uint8(0x20)
				addr := uint16(0xFF00) + uint16(offset)
				bus.Write(addr, 0x99)
				bus.Write(cpu.registers.pc, 0xF0)
				bus.Write(cpu.registers.pc+1, offset)
				cpu.Clock()
				assert.Equal(t, uint8(0x99), cpu.registers.a, "LDH A,(n) should load value from 0xFF00+n")
			})

			// Test LDH (n),A: opcode 0xE0, writes A into 0xFF00+n.
			t.Run("LDH_n_from_A", func(t *testing.T) {
				bus, cpu := setupCPU()
				cpu.registers.a = 0xAB
				offset := uint8(0x30)
				addr := uint16(0xFF00) + uint16(offset)
				bus.Write(cpu.registers.pc, 0xE0)
				bus.Write(cpu.registers.pc+1, offset)
				cpu.Clock()
				assert.Equal(t, uint8(0xAB), bus.Read(addr), "LDH (n),A should store A into memory at 0xFF00+n")
			})
		})
	})

	// Group: Rotation Operations
	t.Run("Rotation", func(t *testing.T) {
		// RLCA (opcode 0x07)
		{
			bus, cpu := setupCPU()
			cpu.registers.a = 0x80
			bus.Write(cpu.registers.pc, 0x07)
			cpu.Clock()
			assert.Equal(t, uint8(0x01), cpu.registers.a, "A should rotate left (RLCA)")
			assert.True(t, cpu.flags.Carry, "Carry flag should be set by RLCA")
		}
		// RRCA (opcode 0x0F)
		{
			bus, cpu := setupCPU()
			cpu.registers.a = 0x01
			bus.Write(cpu.registers.pc, 0x0F)
			cpu.Clock()
			assert.Equal(t, uint8(0x80), cpu.registers.a, "A should rotate right (RRCA)")
			// Depending on implementation, adjust carry flag check as needed.
		}
	})

	// Group: Arithmetic/Logic Instructions
	t.Run("ALU", func(t *testing.T) {
		t.Run("Arithmetic_Logic", func(t *testing.T) {
			// Table-driven tests for various arithmetic/logic instructions.
			// Assumed opcodes:
			// ADD A, B : 0x80, ADC A, B : 0x88, SUB A, B : 0x90,
			// AND A, B : 0xA0, XOR A, B : 0xA8, OR A, B : 0xB0, CP A, B : 0xB8.
			type flags struct {
				zero, carry, halfCarry, subtract bool
			}
			type testCase struct {
				name            string
				opcode          uint8
				initA, initB    uint8
				initCarry       bool  // For ADC; ignored otherwise.
				expectedA       uint8 // Expected value in A (except for CP)
				expectedFlags   flags
				checkAUnchanged bool // For CP: A remains unchanged.
			}
			tests := []testCase{
				// ADD A,B tests
				{"ADD_A_B_simple", 0x80, 1, 2, false, 3, flags{zero: false, carry: false, halfCarry: false, subtract: false}, false},
				{"ADD_A_B_overflow", 0x80, 0xFF, 1, false, 0, flags{zero: true, carry: true, halfCarry: true, subtract: false}, false},
				// ADC A,B tests (with initial carry)
				{"ADC_A_B_simple", 0x88, 1, 2, true, 4, flags{zero: false, carry: false, halfCarry: false, subtract: false}, false},
				{"ADC_A_B_overflow", 0x88, 0xFF, 0, true, 0, flags{zero: true, carry: true, halfCarry: true, subtract: false}, false},
				// SUB A,B tests
				{"SUB_A_B_simple", 0x90, 5, 3, false, 2, flags{zero: false, carry: false, halfCarry: false, subtract: true}, false},
				{"SUB_A_B_zero", 0x90, 3, 3, false, 0, flags{zero: true, carry: false, halfCarry: false, subtract: true}, false},
				// AND A,B test (commonly sets half-carry)
				{"AND_A_B", 0xA0, 0x55, 0xF0, false, 0x50, flags{zero: false, carry: false, halfCarry: true, subtract: false}, false},
				// XOR A,B test
				{"XOR_A_B", 0xA8, 0xFF, 0x0F, false, 0xF0, flags{zero: false, carry: false, halfCarry: false, subtract: false}, false},
				// OR A,B test
				{"OR_A_B", 0xB0, 0x55, 0xAA, false, 0xFF, flags{zero: false, carry: false, halfCarry: false, subtract: false}, false},
				// CP A,B tests (A remains unchanged)
				{"CP_A_B_equal", 0xB8, 3, 3, false, 3, flags{zero: true, carry: false, halfCarry: false, subtract: true}, true},
				{"CP_A_B_diff", 0xB8, 4, 3, false, 4, flags{zero: false, carry: false, halfCarry: false, subtract: true}, true},
			}

			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					bus, cpu := setupCPU()
					cpu.registers.a = tc.initA
					cpu.registers.b = tc.initB
					// For ADC, set the carry flag as needed.
					cpu.flags.Carry = tc.initCarry
					bus.Write(cpu.registers.pc, tc.opcode)
					cpu.Clock()
					if !tc.checkAUnchanged {
						assert.Equal(t, tc.expectedA, cpu.registers.a, tc.name+": A value")
					} else {
						// For CP, A remains unchanged.
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

	// Group: Flow Control Operations
	t.Run("Flow", func(t *testing.T) {
		// Test PUSH BC (opcode 0xC5)
		t.Run("PUSH_BC", func(t *testing.T) {
			bus, cpu := setupCPU()
			cpu.registers.b = 0x12
			cpu.registers.c = 0x34
			cpu.registers.sp = 0xFFFE
			bus.Write(cpu.registers.pc, 0xC5) // PUSH BC opcode
			cpu.Clock()
			// Expect SP decremented by 2.
			assert.Equal(t, uint16(0xFFFC), cpu.registers.sp, "SP should decrease by 2 after PUSH")
			// Assume high byte is stored at SP and low byte at SP+1.
			high := bus.Read(cpu.registers.sp)
			low := bus.Read(cpu.registers.sp + 1)
			assert.Equal(t, uint8(0x12), high, "PUSH_BC: high byte")
			assert.Equal(t, uint8(0x34), low, "PUSH_BC: low byte")
		})

		// Test POP BC (opcode 0xC1)
		t.Run("POP_BC", func(t *testing.T) {
			bus, cpu := setupCPU()
			cpu.registers.sp = 0xFFFC
			// Preload stack with known value: high then low.
			bus.Write(cpu.registers.sp, 0x78)
			bus.Write(cpu.registers.sp+1, 0x9A)
			bus.Write(cpu.registers.pc, 0xC1) // POP BC opcode
			cpu.Clock()
			assert.Equal(t, uint16(0xFFFE), cpu.registers.sp, "SP should increase by 2 after POP")
			assert.Equal(t, uint8(0x78), cpu.registers.b, "POP_BC: register B")
			assert.Equal(t, uint8(0x9A), cpu.registers.c, "POP_BC: register C")
		})

		// Test CALL nn (opcode 0xCD)
		t.Run("CALL_nn", func(t *testing.T) {
			bus, cpu := setupCPU()
			cpu.registers.sp = 0xFFFE
			initPC := cpu.registers.pc
			// CALL target 0x1234 (little-endian: 0x34, 0x12)
			bus.Write(cpu.registers.pc, 0xCD) // CALL nn opcode
			bus.Write(cpu.registers.pc+1, 0x34)
			bus.Write(cpu.registers.pc+2, 0x12)
			cpu.Clock()
			// After CALL, PC should equal target and SP decremented by 2.
			assert.Equal(t, uint16(0x1234), cpu.registers.pc, "CALL should jump to target address")
			assert.Equal(t, uint16(0xFFFC), cpu.registers.sp, "CALL should push return address onto stack")
			// Return address should be initPC + 3.
			retHigh := bus.Read(cpu.registers.sp)
			retLow := bus.Read(cpu.registers.sp + 1)
			expectedRet := initPC + 3
			actualRet := uint16(retHigh)<<8 | uint16(retLow)
			assert.Equal(t, expectedRet, actualRet, "CALL should push correct return address")
		})

		// Test RET (opcode 0xC9)
		t.Run("RET", func(t *testing.T) {
			bus, cpu := setupCPU()
			cpu.registers.sp = 0xFFFC
			// Preload return address: 0x4567.
			bus.Write(cpu.registers.sp, 0x45)
			bus.Write(cpu.registers.sp+1, 0x67)
			bus.Write(cpu.registers.pc, 0xC9) // RET opcode
			cpu.Clock()
			assert.Equal(t, uint16(0x4567), cpu.registers.pc, "RET should set PC to return address")
			assert.Equal(t, uint16(0xFFFE), cpu.registers.sp, "RET should pop return address from stack")
		})

		// Test SBC A, d8 (opcode 0xDE)
		t.Run("SBC_A_d8", func(t *testing.T) {
			// Table-driven tests for SBC A,d8.
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
					bus, cpu := setupCPU()
					cpu.registers.a = tt.initA
					cpu.flags.Carry = tt.initCarry
					bus.Write(cpu.registers.pc, 0xDE) // SBC A, d8 opcode
					bus.Write(cpu.registers.pc+1, tt.immediate)
					cpu.Clock()
					assert.Equal(t, tt.expectedA, cpu.registers.a, "SBC A,d8 result mismatch")
					assert.Equal(t, tt.expectedZero, cpu.flags.Zero, "SBC A,d8 zero flag mismatch")
				})
			}
		})

		// Test JP nn (opcode 0xC3)
		t.Run("JP_nn", func(t *testing.T) {
			bus, cpu := setupCPU()
			// Jump to 0x4321 (little-endian: 0x21, 0x43).
			bus.Write(cpu.registers.pc, 0xC3) // JP nn opcode
			bus.Write(cpu.registers.pc+1, 0x21)
			bus.Write(cpu.registers.pc+2, 0x43)
			cpu.Clock()
			assert.Equal(t, uint16(0x4321), cpu.registers.pc, "JP should jump to immediate address")
		})

		// Test RST (opcode 0xDF for RST 18h)
		t.Run("RST", func(t *testing.T) {
			bus, cpu := setupCPU()
			cpu.registers.sp = 0xFFFE
			initPC := cpu.registers.pc
			bus.Write(cpu.registers.pc, 0xDF)
			cpu.Clock()
			assert.Equal(t, uint16(0x0018), cpu.registers.pc, "RST should set PC to fixed vector 0x0018")
			// Verify that return address (initPC+1) is pushed.
			assert.Equal(t, uint16(0xFFFC), cpu.registers.sp, "RST should push return address onto stack")
			retHigh := bus.Read(cpu.registers.sp)
			retLow := bus.Read(cpu.registers.sp + 1)
			expectedRet := initPC + 1
			actualRet := uint16(retHigh)<<8 | uint16(retLow)
			assert.Equal(t, expectedRet, actualRet, "RST should push correct return address")
		})
	})
}
