package io

// --- End stubs ---

// func TestRegistersTableReadWrite(t *testing.T) {
// 	fb := &FakeBus{}
// 	regs := NewRegisters(nil, fb, nil)

// 	// Define test cases per predefined register address.
// 	// For read/write registers we expect the read value to match the written one (or its converted form).
// 	// For read-only registers (like LY at 0x44) we expect a panic.
// 	tests := []struct {
// 		name         string
// 		addr         uint16
// 		writeVal     uint8
// 		expectedRead uint8 // expected read result after writing
// 		shouldPanic  bool
// 	}{
// 		{"Joypad (0xFF00)", 0x00, 0xAB, 0xAB, false},
// 		{"Serial byte 0 (0xFF01)", 0x01, 0x11, 0x11, false},
// 		{"Serial byte 1 (0xFF02)", 0x02, 0x22, 0x22, false},
// 		{"Timer DIV (0xFF04)", 0x04, 0x33, 0x33, false},
// 		{"Timer TIMA (0xFF05)", 0x05, 0x44, 0x44, false},
// 		{"Timer TMA (0xFF06)", 0x06, 0x55, 0x55, false},
// 		{"Timer TAC (0xFF07)", 0x07, 0x66, 0x66, false},
// 		{"Interrupt Flag (0xFF0F)", 0x0F, 0x77, 0x77, false},
// 		{"Audio (0xFF10)", 0x10, 0x88, 0x88, false},
// 		{"Wave Pattern (0xFF30)", 0x30, 0x99, 0x99, false},
// 		{"LCD Control (0xFF40)", 0x40, 0xAA, 0xAA, false},
// 		{"LCD Status (0xFF41)", 0x41, 0xBB, 0xBB, false},
// 		{"ScrollY (0xFF42)", 0x42, 0xCC, 0xCC, false},
// 		{"ScrollX (0xFF43)", 0x43, 0xDD, 0xDD, false},
// 		// LY (0xFF44) is read-only
// 		{"LY (0xFF44) read-only", 0x44, 0xFF, 0, true},
// 		{"LY Compare (0xFF45)", 0x45, 0xEE, 0xEE, false},
// 		// DMA (0xFF46) does a DMA transfer, so its write does not store a value.
// 		{"DMA (0xFF46)", 0x46, 0x10, 0, false},
// 		{"Palette Data (0xFF47)", 0x47, 0xF1, 0xF1, false},
// 		{"OBJ Palette 0 (0xFF48)", 0x48, 0xF2, 0xF2, false},
// 		{"OBJ Palette 1 (0xFF49)", 0x49, 0xF3, 0xF3, false},
// 		{"Window Y (0xFF4A)", 0x4A, 0xF4, 0xF4, false},
// 		{"Window X (0xFF4B)", 0x4B, 0xF5, 0xF5, false},
// 		// DisableBootROM (0xFF50) as boolean (nonzero → 1)
// 		{"DisableBootROM (0xFF50)", 0x50, 0x01, 0x01, false},
// 		// VRAMBank1 (0xFF4F) as boolean
// 		{"VRAMBank1 (0xFF4F)", 0x4F, 0x01, 0x01, false},
// 		// GBCPaletteData (0xFF68) test first byte
// 		{"GBCPaletteData (0xFF68)", 0x68, 0xAA, 0xAA, false},
// 		// VRAMDMA (0xFF51) test first byte
// 		{"VRAMDMA (0xFF51)", 0x51, 0xBB, 0xBB, false},
// 		// WRAMBank1 (0xFF70) as boolean
// 		{"WRAMBank1 (0xFF70)", 0x70, 0x01, 0x01, false},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			if tc.shouldPanic {
// 				defer func() {
// 					if r := recover(); r == nil {
// 						t.Errorf("Expected panic for register write at 0x%X", tc.addr)
// 					}
// 				}()
// 				regs.Write(tc.addr, tc.writeVal)
// 			} else {
// 				// Write the test value.
// 				regs.Write(tc.addr, tc.writeVal)
// 				// Read back.
// 				got := regs.Read(tc.addr)
// 				if got != tc.expectedRead {
// 					t.Errorf("addr 0x%X: expected 0x%02X, got 0x%02X", tc.addr, tc.expectedRead, got)
// 				}
// 			}
// 		})
// 	}
// }
