// Package I2cKbd creates an interface to the keyboard of the
// picocalc

//go:build tinygo

package lib

import (
	"fmt"
	"machine"
)

// IMPORTANT:
//
// 1. The PicoCalc must be powered on for the i2c keyboard chip to be active.
// I wasted a bit of time discovering this!
//
// 2. Do not use batteries in the PicoCalc while it's pluggen in via USB mini and turned on
// becuase an electrical path is opened that causes the 18650 batteries to be charged
// beyond 4.2 volts.  Hopefully they fix this hardware flaw.
var i2cKbdAddr uint16 = 0x1F

const i2cGetKey = 0x09

const (
	ALT_KEY       byte = 0xA1
	BACKSPACE_KEY      = 0x08
	CTRL_KEY      byte = 0xA5
	DEL_KEY       byte = 0xd4
	END_KEY       byte = 0xd5
	ESC_KEY       byte = 0xb1
	F1_KEY        byte = 0x81
	F2_KEY        byte = 0x82
	F3_KEY        byte = 0x83
	F4_KEY        byte = 0x84
	F5_KEY        byte = 0x85
	F6_KEY        byte = 0x86
	F7_KEY        byte = 0x87
	F8_KEY        byte = 0x88
	F9_KEY        byte = 0x89
	F10_KEY       byte = 0x90 // odd it's not 0x8A
	HOME_KEY      byte = 0xd2
	INS_KEY       byte = 0xd1

	LEFT_KEY  byte = 0xb4
	RIGHT_KEY byte = 0xb7
	UP_KEY    byte = 0xb5
	DOWN_KEY  byte = 0xb6

	/*KEY_SUP    // 281
	KEY_SDOWN  // 282
	KEY_SLEFT  // 283
	KEY_SRIGHT // 284
	KEY_SHOME  // 285
	KEY_SEND   // 286
	// editing
	KEY_CUT   // 287
	KEY_COPY  // 288
	KEY_PASTE // 289
	KEY_SAVE  // 290
	KEY_QUIT  // 291
	KEY_HELP  // 292*/

)

type I2CKbd struct {
	i2c      *machine.I2C
	write    []byte
	read     []byte
	AltDown  bool
	CtrlDown bool
}

// Init initialized the i2c driver.  It may be necessary to add the ability to
// provided an i2c driver if the bus is shared (I don't believe it is currently).
func (kbd *I2CKbd) Init() error {
	kbd.write = make([]byte, 1)
	kbd.write[0] = i2cGetKey
	kbd.read = make([]byte, 2)
	kbd.i2c = machine.I2C1
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
		return kbd.ifNoModifiers(27)
	/*case 'C':
		if kbd.AltDown {
			return key.KEY_COPY
		}
	case 'H':
		if kbd.AltDown {
			return key.KEY_HELP
		}
	case 'Q':
		if kbd.AltDown {
			return key.KEY_QUIT
		}
	case 'S':
		if kbd.AltDown {
			return key.KEY_SAVE
		}
	case 'V':
		if kbd.AltDown {
			return key.KEY_PASTE
		}
	case 'X':
		if kbd.AltDown {
			return key.KEY_CUT
		}
	case 'c':
		if kbd.CtrlDown {
			return key.KEY_BREAK
		}*/
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
	return 0
}

// Called when a key is released.  We mostly don't care outside of modifier keys
func (kbd *I2CKbd) keyUp() byte {
	switch kbd.read[1] {
	case ALT_KEY:
		kbd.AltDown = false
	case CTRL_KEY:
		kbd.CtrlDown = false
	}
	return 0
}

// Covers the common path where we only want to report a key
// if no other modifiers are held down.
func (kbd *I2CKbd) ifNoModifiers(k byte) byte {
    if kbd.CtrlDown || kbd.AltDown {
        return 0
    }
    return k
}

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
