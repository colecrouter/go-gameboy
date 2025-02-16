package registers

type TransferRate uint8

const (
	TransferRateNormal TransferRate = iota // 8192Hz
	// TODO GBC
)

type SerialDevice interface {
	Transfer(uint8) uint8
}

type SerialTransfer struct {
	EnableTransfer bool
	TransferRate   TransferRate
	Master         bool

	connected SerialDevice
	interrupt *Interrupt
	value     uint8
} // 0xFF01-0xFF02

func NewSerialTransfer(interupt *Interrupt) *SerialTransfer {
	return &SerialTransfer{
		interrupt: interupt,
	}
}

func (s *SerialTransfer) Read(addr uint16) uint8 {
	switch addr {
	case 0x00:
		return s.value
	case 0x01:
		var val uint8
		if s.EnableTransfer {
			val |= 1 << 7
		}
		switch s.TransferRate {
		case TransferRateNormal:
			// no bits set
		}
		if s.Master {
			val |= 1 << 0
		}
		return val
	}
	panic("Invalid serial transfer register")
}

func (s *SerialTransfer) Write(addr uint16, val uint8) {
	switch addr {
	case 0x00:
		s.value = val
	case 0x01:
		s.EnableTransfer = val&0x80 != 0
		s.TransferRate = TransferRate(val & 0x1)
		s.Master = val&0x1 != 0

		// Transfer subroutine
		if s.EnableTransfer {
			// go s.transfer()
			s.transfer()
		}
	default:
		panic("Invalid serial transfer register")
	}
}

func (s *SerialTransfer) Connect(device SerialDevice) {
	s.connected = device
}

func (s *SerialTransfer) transfer() {
	if s.connected != nil {
		s.value = s.connected.Transfer(s.value)
	}

	// Clear transfer bit
	s.EnableTransfer = false

	// Trigger interupt
	if s.interrupt != nil {
		s.interrupt.Serial = true
	}
}
