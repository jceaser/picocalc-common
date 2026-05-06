//go:build tinygo

// Package I2cKbd creates an interface to the keyboard of the
// picocalc

package i2ckbd

import (
	"fmt"
	"machine"
)


// ***************************************************************************80
// Battery

func (kbd *I2CKbd) ReadBattery() int16 {
    msg := []byte{0x0B} // Command to send
	//kbd.write[0] = 0x0B
	
    err := kbd.i2c.Tx(i2cKbdAddr, msg, kbd.read)
	if err != nil {
		return -1
	}
	
	if (kbd.read[0] == 0) && (kbd.read[1] == 0) {
		return 0
	}

    return int16(kbd.read[1])
}
