//go:build tinygo

// ********************************************************************************************* 100
/*
Package I2cKbd creates an interface to the keyboard of the
picocalc

lots of details can be found at:
https://pip-assets.raspberrypi.com/categories/814-rp2040/documents/RP-008371-DS-1-rp2040-datasheet.pdf
*/

package i2ckbd

import (
	"fmt"
	"machine"
)

// IMPORTANT:
//
// 1. The PicoCalc must be powered on for the i2c keyboard chip to be active.
// I [mattwach] wasted a bit of time discovering this!
//
// 2. Do not use batteries in the PicoCalc while it's pluggen in via USB mini and turned on
// becuase an electrical path is opened that causes the 18650 batteries to be charged
// beyond 4.2 volts.  Hopefully they fix this hardware flaw.
var i2cKbdAddr uint16 = 0x1F

const i2cGetKey = 0x09

const (
	BACKSPACE_KEY 	byte = 0x08
	TAB_KEY		  	byte = 0x09
	ENTER_KEY		byte = 0x0a

	ESC        		byte = 0x1b
	SPACE_KEY		byte = 0x20

	// ! - @ 	is 0x21 to 0x40
	// A - Z 	is 0x41 to 0x5a
	// [ - ~ 	is 0x5b to 0x60
	// a - z	is 0x61 to 0x7a
	// { - ~ 	is 0x7b to 0x7e
	// no 0x7f or 0x80

	F1_KEY        	byte = 0x81 // held
	F2_KEY        	byte = 0x82 // held
	F3_KEY        	byte = 0x83 // held
	F4_KEY        	byte = 0x84 // held
	F5_KEY        	byte = 0x85 // held
	F6_KEY        	byte = 0x86 // held
	F7_KEY        	byte = 0x87 // held
	F8_KEY        	byte = 0x88 // held
	F9_KEY        	byte = 0x89 // held
	F10_KEY       	byte = 0x90 // held - odd it's not 0x8A

	ALT_KEY       	byte = 0xa1 // held
	LEFT_SHIFT_KEY 	byte = 0xa2 // held
	RIGHT_SHIFT_KEY byte = 0xa3 // held
	CTRL_KEY      	byte = 0xa5 // held

	ESC_KEY			byte = 0xb1 // held - note it is not 0x1b

	LEFT_KEY		byte = 0xb4
	UP_KEY			byte = 0xb5
	DOWN_KEY   	  	byte = 0xb6
	RIGHT_KEY 	  	byte = 0xb7

	CAPS_LOCK_KEY	byte = 0xc1 // held

	BREAK_KEY		byte = 0xd0
	INS_KEY       	byte = 0xd1
	HOME_KEY      	byte = 0xd2 // 210
	DEL_KEY       	byte = 0xd4
	END_KEY       	byte = 0xd5
)

type I2CKbd struct {
	i2c      *machine.I2C
	write    []byte
	read     []byte
	AltDown  bool
	CtrlDown bool
	LastKey  KeyDetail
}

// ***************************************************************************80

// Init initialized the i2c driver.  It may be necessary to add the ability to
// provided an i2c driver if the bus is shared (I don't believe it is currently).
func (kbd *I2CKbd) Init() error {
	kbd.write = make([]byte, 1)
	kbd.write[0] = i2cGetKey
	kbd.read = make([]byte, 2)
	kbd.i2c = machine.I2C1
	kbd.LastKey = KeyDetail{}
	return kbd.i2c.Configure(machine.I2CConfig{
		SCL: machine.GP7,
		SDA: machine.GP6,
	})
}

// GetChar returns a keppress using the key driver codes.  Returns zero if
// nothing was pressed (which will be most of the time).
//
// You need to call this often.  Calling it in a gorouting seems like a good
// plan, but should be done after the basics are fully sorted.
func (kbd *I2CKbd) GetChar() (byte, error) {
	err := kbd.i2c.Tx(i2cKbdAddr, kbd.write, kbd.read)
	if err != nil {
		return 0, err
	}
	if (kbd.read[0] == 0) && (kbd.read[1] == 0) {
		return 0, nil
	}
	switch kbd.read[0] {
	case 0x01:
		return kbd.keyDown(), nil
	case 0x02:
		return kbd.keyHeld(), nil
	case 0x03:
		return kbd.keyUp(), nil
	default:
		return 0, fmt.Errorf("unknown key response: %v", kbd.read[0])
	}
}

// called when a key is depressed
func (kbd *I2CKbd) keyDown() byte {
	k := kbd.read[1]
	switch k {
	case ALT_KEY:
		kbd.AltDown = true
		return 0
	case CTRL_KEY:
		kbd.CtrlDown = true
		return 0
	case F1_KEY:
		return kbd.ifNoModifiers(F1_KEY)
	case F2_KEY:
		return kbd.ifNoModifiers(F2_KEY)
	case F3_KEY:
		return kbd.ifNoModifiers(F3_KEY)
	case F4_KEY:
		return kbd.ifNoModifiers(F4_KEY)
	case F5_KEY:
		return kbd.ifNoModifiers(F5_KEY)
	case F6_KEY:
		return kbd.ifNoModifiers(F6_KEY)
	case F7_KEY:
		return kbd.ifNoModifiers(F7_KEY)
	case F8_KEY:
		return kbd.ifNoModifiers(F8_KEY)
	case F9_KEY:
		return kbd.ifNoModifiers(F9_KEY)
	case F10_KEY:
		return kbd.ifNoModifiers(F10_KEY)
	case LEFT_KEY:
		/*if kbd.AltDown {
			return LEFT_KEY
		}*/
		return kbd.ifNoModifiers(LEFT_KEY)
	case RIGHT_KEY:
		/*if kbd.AltDown {
			return key.KEY_SRIGHT
		}*/
		return kbd.ifNoModifiers(RIGHT_KEY)
	case UP_KEY:
		/*if kbd.AltDown {
			return key.KEY_SUP
		}
		if kbd.CtrlDown {
			return key.KEY_PAGEUP
		}*/
		return UP_KEY
	case DOWN_KEY:
		/*if kbd.AltDown {
			return key.KEY_SDOWN
		}
		if kbd.CtrlDown {
			return key.KEY_PAGEDOWN
		}*/
		return DOWN_KEY
	case BACKSPACE_KEY:
		return kbd.ifNoModifiers(BACKSPACE_KEY)
	case DEL_KEY:
		return kbd.ifNoModifiers(DEL_KEY)
	case INS_KEY:
		return kbd.ifNoModifiers(INS_KEY)
	case END_KEY:
		if kbd.AltDown {
			return END_KEY //key.KEY_SEND
		}
		return kbd.ifNoModifiers(END_KEY)
	case HOME_KEY:
		/*if kbd.AltDown {
			return key.KEY_SHOME
		}*/
		return kbd.ifNoModifiers(HOME_KEY)
	case ESC_KEY:
		return kbd.ifNoModifiers(ESC)
	}
	if k < 0x80 {
		return kbd.ifNoModifiers(k)
	}
	return 0
}

// Sometimes called when a key is held.  Usually just for modifier keys.
func (kbd *I2CKbd) keyHeld() byte {
	switch kbd.read[1] {
	case ALT_KEY:
		// likely not needed, but doesn't hurt anything either
		kbd.AltDown = true
	case CTRL_KEY:
		// likely not needed, but doesn't hurt anything either
		kbd.CtrlDown = true
	}
	return 0x02
}

// Called when a key is released.  We mostly don't care outside of modifier keys
func (kbd *I2CKbd) keyUp() byte {
	switch kbd.read[1] {
	case ALT_KEY:
		kbd.AltDown = false
	case CTRL_KEY:
		kbd.CtrlDown = false
	}
	return 0x03
}

// Covers the common path where we only want to report a key
// if no other modifiers are held down.
func (kbd *I2CKbd) ifNoModifiers(k byte) byte {
    if kbd.CtrlDown || kbd.AltDown {
        return 0
    }
    return k
}
