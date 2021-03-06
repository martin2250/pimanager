package gammapwm

import (
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

type GammaPWM struct {
	Bus     string
	Address int
	Value   [8]byte
}

func (d *GammaPWM) Update() error {
	bus, err := i2creg.Open(d.Bus)
	if err != nil {
		return err
	}
	defer bus.Close()

	dev := i2c.Dev{Bus: bus, Addr: uint16(d.Address)}

	buffer := make([]byte, 9)

	copy(buffer[1:], d.Value[:])

	_, err = dev.Write(buffer)
	return err
}

func (d *GammaPWM) Init() error {
	bus, err := i2creg.Open(d.Bus)
	if err != nil {
		return err
	}
	defer bus.Close()

	dev := i2c.Dev{Bus: bus, Addr: uint16(d.Address)}

	buffer := make([]byte, 8)

	if err = dev.Tx([]byte{0}, buffer); err != nil {
		return err
	}

	copy(d.Value[:], buffer[:])

	return err
}
