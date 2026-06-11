//go:build tinygo

package i2ckbd

import (
	"fmt"
)

/*

1. how do I preserve action keys alt, ctl, shift between multiple key presses, need a global setting for them.

*/

type EventType byte
const (
	EVENT_NONE EventType = iota
	EVENT_DOWN
	EVENT_HELD
	EVENT_UP
)

// ***************************************************************************80

type KeyDetail struct {
	Key byte
	Action EventType
	modifiers byte // ? alt control shift
}

/*
This version should always return the key that was pressed along with the flags
*/
func (kbd *I2CKbd) Event() (KeyDetail, error) {
	key := KeyDetail{}
	err := kbd.i2c.Tx(i2cKbdAddr, kbd.write, kbd.read)
	if err != nil {
		return key, err
	}

	if (kbd.read[0] == 0) && (kbd.read[1] == 0) {
		return key, nil
	}

	key.Action = EventType(kbd.read[0])
	key.Key = kbd.read[1]

	//fmt.Printf("Event() = B: %s.\n", key)
	key.Copy(kbd.LastKey)
	//fmt.Printf("Event() = A: %s.\n", key)

	switch key.Action {
	case EVENT_DOWN:
		if key.Modifier() {
			key.Set()
			kbd.LastKey.SetFromKey(key)
		}
		fmt.Printf("Event() = R: %s.\n", key)
		fmt.Printf("Event() = L: %s.\n", kbd.LastKey)
		return key, nil
	case EVENT_HELD:
		return key, nil
	case EVENT_UP:
		if key.Modifier() {
			key.Clear()
			kbd.LastKey.ClearFromKey(key)
		}
		return key, nil
	default:
		return key, fmt.Errorf("Error unknown key response: %v", kbd.read[0])
	}
}

func (k *KeyDetail) Copy(existing KeyDetail) {
	if existing.Control() {
		k.SetControl()
	} else {
		k.ClearControl()
	}
	if existing.RightShift() {
		k.SetRightShift()
	} else {
		k.ClearRightShift()
	}
	if existing.LeftShift() {
		k.SetLeftShift()
	} else {
		k.ClearLeftShift()
	}
	if existing.Alt() {
		k.SetAlt()
	} else {
		k.ClearAlt()
	}
}

func (k KeyDetail) Modifier() bool {
	return k.Key == ALT_KEY || k.Key == LEFT_SHIFT_KEY || k.Key == RIGHT_SHIFT_KEY || k.Key == CTRL_KEY
}

// *************************************
// Setters

func (k *KeyDetail) SetControl() {
	k.modifiers = k.modifiers | 0x08
}

func (k *KeyDetail) SetRightShift() {
	k.modifiers = k.modifiers | 0x04
}

func (k *KeyDetail) SetLeftShift() {
	k.modifiers = k.modifiers | 0x02
}

func (k *KeyDetail) SetAlt() {
	k.modifiers = k.modifiers | 0x01
}

// *************************************
// Clears

func (k *KeyDetail) ClearControl() {
	k.modifiers = k.modifiers & 0xf7 // 1111 0111
}

func (k *KeyDetail) ClearRightShift() {
	k.modifiers = k.modifiers & 0xfb // 1111 1011
}

func (k *KeyDetail) ClearLeftShift() {
	k.modifiers = k.modifiers & 0xfd // 1111 1101
}

func (k *KeyDetail) ClearAlt() {
	k.modifiers = k.modifiers & 0xfe // 1111 1110
}

// *************************************
// Getters

func (k KeyDetail) Control() bool {
	return k.modifiers & 0x08 > 0
}

func (k KeyDetail) RightShift() bool {
	return k.modifiers & 0x04 > 0
}

func (k KeyDetail) LeftShift() bool {
	return k.modifiers & 0x02 > 0
}

func (k KeyDetail) Alt() bool {
	return k.modifiers & 0x01 > 0
}

// *************************************

func (k *KeyDetail) ClearFromKey(other KeyDetail) {
	if other.Key == CTRL_KEY {
		k.ClearControl()
	} else if other.Key == LEFT_SHIFT_KEY {
		k.ClearLeftShift()
	} else if other.Key == RIGHT_SHIFT_KEY {
		k.ClearRightShift()
	} else if other.Key == ALT_KEY {
		k.ClearAlt()
	}
}

func (k *KeyDetail) Clear() {
	if k.Key == CTRL_KEY {
		k.ClearControl()
	} else if k.Key == LEFT_SHIFT_KEY {
		k.ClearLeftShift()
	} else if k.Key == RIGHT_SHIFT_KEY {
		k.ClearRightShift()
	} else if k.Key == ALT_KEY {
		k.ClearAlt()
	}
}

func (k *KeyDetail) SetFromKey(other KeyDetail) {
	if other.Key == CTRL_KEY {
		k.SetControl()
	} else if other.Key == LEFT_SHIFT_KEY {
		k.SetLeftShift()
	} else if other.Key == RIGHT_SHIFT_KEY {
		k.SetRightShift()
	} else if other.Key == ALT_KEY {
		k.SetAlt()
	}
}

func (k *KeyDetail) Set() {
	if k.Key == CTRL_KEY {
		k.SetControl()
	} else if k.Key == LEFT_SHIFT_KEY {
		k.SetLeftShift()
	} else if k.Key == RIGHT_SHIFT_KEY {
		k.SetRightShift()
	} else if k.Key == ALT_KEY {
		k.SetAlt()
	}
}

func (k KeyDetail) String() string {
	alt := ""
	if k.Alt() {
		alt = "Alt"
	}
	ctl := ""
	if k.Control() {
		ctl = "Control"
	}
	shift := ""
	if k.LeftShift() || k.RightShift() {
		shift = "Shift"
	}
	out := fmt.Sprintf("key: 0x%x, Action: 0x%x, Flags: %d [%s %s %s]",
		k.Key,
		k.Action,
		k.modifiers,
		shift,
		ctl,
		alt,
		)
	return out
}
