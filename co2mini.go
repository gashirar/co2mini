package co2mini

import (
	"fmt"
	"time"
	"github.com/zserge/hid"
	"errors"
)

var (
	key = []byte{0x86, 0x41, 0xc9, 0xa8, 0x7f, 0x41, 0x3c, 0xac}
)

type Co2mini struct {
	Co2Ch  chan int
	TempCh chan float64
	device hid.Device
}

func (c *Co2mini) Connect() error {
	const vendorId = "04d9:a052:0100:00"

	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		id := fmt.Sprintf("%04x:%04x:%04x:%02x", info.Vendor, info.Product, info.Revision, info.Interface)
		if id != vendorId {
			return
		}
		c.device = device
	})
	if c.device == nil {
		return errors.New("Device not found.")
	}
	c.Co2Ch = make(chan int)
	c.TempCh = make(chan float64)
	return nil
}

func (c *Co2mini) Start() error {
	const co2op = 0x50
	const tempop = 0x42

	if err := c.device.Open(); err != nil {
		return err
	}
	defer c.device.Close()

	if err := c.device.SetReport(0, key); err != nil {
		return err
	}

	for {
		if buf, err := c.device.Read(-1, 1*time.Second); err == nil {
			dec := decrypt(buf, key)
			if len(dec) == 0 {
				continue
			}
			val := int(dec[1])<<8 | int(dec[2])
			if dec[0] == co2op {
				c.Co2Ch <- val
			}
			if dec[0] == tempop {
				c.TempCh <- float64(val)/16.0 - 273.15
			}
		}
	}
}

func decrypt(b, key []byte) []byte {
	if len(b) != 8 {
		return b
	}
	phase1 := shuffle(b)
	phase2 := xor(phase1, key)
	phase3 := shift(phase2)
	ctmp := offset()
	result := calc(phase3, ctmp)
	return result
}

func calc(b, ctmp []byte) []byte {
	res := make([]byte, 8)
	for i := range b {
		res[i] = (0xFF + b[i] - ctmp[i] + 0x01) & 0xFF
	}
	return res
}

func offset() []byte {
	offset := []byte{0x48, 0x74, 0x65, 0x6D, 0x70, 0x39, 0x39, 0x65} //"Htemp99e"
	res := make([]byte, 8)
	for i := range offset {
		res[i] = (offset[i]>>4 | offset[i]<<4) & 0xFF
	}
	return res
}

func shift(b []byte) []byte {
	res := make([]byte, 8)
	for i := range b {
		res[i] = (b[i]>>3 | b[(i-1+8)%8]<<5) & 0xFF
	}
	return res
}

func xor(b, key []byte) []byte {
	res := make([]byte, 8)
	for i := range b {
		res[i] = b[i] ^ key[i]
	}
	return res
}

func shuffle(b []byte) []byte {
	assignNum := []int{2, 4, 0, 7, 1, 6, 5, 3}
	res := make([]byte, 8)
	for i, v := range assignNum {
		res[i] = b[v]
	}
	return res
}
