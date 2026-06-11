//go:build tinygo

/************************************************************************************************100
This is a test app which pulls in the different packages and runs a nominal execution of the code
just to make sure everything is working at the most basic of levels. The code can also be used as an
example.

Compile with:
	 tinygo build -o test.uf2 -target=pico
	 tinygo build -o test2.uf2 -target=pico2
***************************************************************************************************/

package main

import (
	"fmt"
	"image/color"
	"machine"
	"math"
	"time"

	"device/rp"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"

	"github.com/jceaser/picocalc-common/i2ckbd"
	"github.com/jceaser/picocalc-common/ili948x"
	"github.com/jceaser/picocalc-common/sound"
	"github.com/jceaser/picocalc-common/util"
)

const (
	FREQ_LOW = 20.0
	FREQ_HIGH = 10_000.0 // I personally can not hear 20_000 on the PicoCalc
	FREQ_DIF = FREQ_HIGH / FREQ_LOW
	PICOCALC_BL_MAGIC = 0xe98cc638 //id to send message to UF2 loader
)

var (
	black = color.RGBA{0, 0, 0, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	white = color.RGBA{255, 255, 255, 255}
	font = &freemono.Regular12pt7b // Regular18pt7b
	lastDebug = ""
)

// ***************************************************************************80

func maximum(one, two int16) int16 {
	if one < two {
		return two
	}
	return one
}

func minimum (one, two int16) int16 {
	if one < two {
		return one
	}
	return two
}

func clearDebug(lcd *ili948x.Ili948x, row int16) {
	column := row * 15 + 15
	if len(lastDebug) > 0 {
		tinyfont.WriteLineRotated(lcd, font, column, 315, lastDebug, black, tinyfont.ROTATION_270)
		lastDebug = ""
	}
}

func debug(lcd *ili948x.Ili948x, row int16, label string, value any) {
	column := row * 15 + 15
	if len(lastDebug) > 0 {
		tinyfont.WriteLineRotated(lcd, font, column, 315, lastDebug, black, tinyfont.ROTATION_270)
	}
	display := fmt.Sprintf("%s = %v", label, value)
	/*if lastDebug == display {
		return
	}*/
	tinyfont.WriteLineRotated(lcd, font, column, 315, display, white, tinyfont.ROTATION_270)
	lastDebug = display
}

// Draw an x shaped pointer with pos at the center
func drawPointer(display *ili948x.Ili948x, pos ili948x.Point, color ili948x.RGB565) {
	p1 := ili948x.Point{X: pos.X-10, Y: pos.Y-10}
	p2 := ili948x.Point{X: pos.X+10, Y: pos.Y+10}
	display.DrawLine(p1, p2, color)

	p3 := ili948x.Point{X: pos.X+10, Y: pos.Y-10}
	p4 := ili948x.Point{X: pos.X-10, Y: pos.Y+10}
	display.DrawLine(p3, p4, color)
}

// remember, 320x320 screen size
func main() {
	if util.Version() == 0.2 {
		machine.Watchdog.Start()
		display := ili948x.InitDisplay()

		var keyboard i2ckbd.I2CKbd
		if err := keyboard.Init(); err != nil {
			fmt.Println(err)
		}
		// test the screen
		display.DrawHLine(0, 320, 1, ili948x.BLUE) // top
		display.DrawVLine(319, 0, 320, ili948x.BLUE) // right
		display.DrawHLine(0, 320, 319, ili948x.BLUE) // bottom
		display.DrawVLine(1, 0, 320, ili948x.BLUE) // left

		// play some sounds
		snd := sound.Initialize()
		sound.Play(sound.SND_BEEP)
		sound.Play(sound.SND_TAB_SWITCH)
		sound.Play(sound.SND_ERROR)

		// draw some more to confirm that sound returned control
		display.DrawHLine(1, 319, 159, ili948x.RED)
		display.DrawVLine(159, 1, 319, ili948x.GREEN)

		// do something interactive to test control is fluid
		pos := ili948x.Point{X: int16(320/4), Y: int16(320/4)}
		last_pos := ili948x.Point{X: pos.X, Y: pos.Y}
		redraw := true
		tick := 0
		fmt.Print("\033[2J") // clear screen
		for {
			tick++
			k, err := keyboard.Event()
			if err != nil {
				fmt.Printf("Error: Event=%v\n", err)
				sound.Play(sound.SND_ERROR)
			}

			// ignore held modifiers
			if k.Action == i2ckbd.EVENT_HELD && k.Modifier() {
				continue
			}

			// ignore down events for everything but arrows
			if k.Action == i2ckbd.EVENT_DOWN {
				if k.Key < i2ckbd.LEFT_KEY || i2ckbd.RIGHT_KEY < k.Key {
					continue
				}
			}

			fmt.Print("\033[s") // save
			fmt.Print("\033[0;0H") // set position
			fmt.Print("\033[2K") // clear
			fmt.Printf("event is %s\n", k)
			fmt.Print("\033[u") // restore

			switch k.Key {
			case 0:
				//unknown, no action
				if k.Key != 0 {
					fmt.Printf("k=%d\n", k.Key)
					debug(display, int16(k.Key), "k=", k.Key)
				}
				continue

			case 0x20:
				if k.Control() {
					fmt.Println("space with control")
				}

			case i2ckbd.ESC_KEY, i2ckbd.ESC:
				clearDebug(display, 0)
				clearDebug(display, 1)
				clearDebug(display, 2)
				clearDebug(display, 3)
				snd.Play(sound.SND_BEEP)

			// do pre-made sounds
			case i2ckbd.F1_KEY:
				snd.Play(sound.SND_ERROR)
			case i2ckbd.F2_KEY:
				snd.Play(sound.SND_BEEP)
			case i2ckbd.F3_KEY:
				snd.Play(sound.SND_TAB_SWITCH)

			// try safely with left speaker
			case i2ckbd.F4_KEY:
				snd.PlayFrequencyLeft(2600, 2000 * time.Millisecond)
			case i2ckbd.F5_KEY:
				snd.PlayFrequencyLeft(32, 2000 * time.Millisecond)

			// send safely to both left and right
			case i2ckbd.F6_KEY:
				snd.PlayFrequencyLeft(2600, 2000 * time.Millisecond)
				snd.PlayFrequencyRight(1200, 2000 * time.Millisecond)

			case i2ckbd.F10_KEY:
				//special code to ask watchdog to give control back to UF2 Loader
				rp.WATCHDOG.SCRATCH1.Set(1)
				rp.WATCHDOG.SCRATCH0.Set(PICOCALC_BL_MAGIC)
				fmt.Printf("🚀 - scratch pad, [%v][%v]\n",
					rp.WATCHDOG.SCRATCH0.Get(),
					rp.WATCHDOG.SCRATCH1.Get())
				machine.CPUReset()

			case i2ckbd.TAB_KEY:
				snd.Play(sound.SND_TAB_SWITCH)

			// try each channel directly
			case i2ckbd.HOME_KEY:
				sound.SetPwmLevelLeft(2600, 200 * time.Millisecond)
			case i2ckbd.END_KEY:
				sound.SetPwmLevelRight(3000, 200 * time.Millisecond)

			case i2ckbd.LEFT_KEY:
				pos.X = maximum(0, pos.X-1)
				redraw = true
			case i2ckbd.RIGHT_KEY:
				pos.X = minimum(319, pos.X+1)
				redraw = true
			case i2ckbd.UP_KEY:
				pos.Y = maximum(0, pos.Y-1)
				redraw = true
			case i2ckbd.DOWN_KEY:
				pos.Y = minimum(319, pos.Y+1)
				redraw = true

			// respond to letter a and b
			/*case 0x61:
				snd.PlayFrequencyLeft(20, 1 * time.Second)
			case 0x62:
				snd.PlayFrequencyRight(788, 1 * time.Second)*/
			case 0x78: //x
				snd.PlayFrequency(5975, 1 * time.Second)
			case 0x79: //y
				snd.PlayFrequency(5980, 1 * time.Second)
			case 0x7a: //z - 6000 is to far
				snd.PlayFrequency(5990, 1 * time.Second)

			default:
				//distribute frequencies between keys to test range
				if 0x61 <= k.Key && k.Key <= 0x7a {
					// a letter z
					exp := float64(k.Key-0x61) / (0x7a - 0x61)
					freq := int(FREQ_LOW * math.Pow(FREQ_DIF, exp))
					fmt.Printf("freq: %d\n", freq)
					snd.PlayFrequency(freq, 1 * time.Second)
				} else {
					fmt.Printf("event= %s\n", k)
				}
			}

			// bounding box error
			if pos.Y == 0 || pos.X == 0 || pos.Y == 319 || pos.X == 319 {
				snd.Play(sound.SND_ERROR)
			}

			// draw pointer if it moved
			if redraw {
				drawPointer(display, last_pos, ili948x.BLACK)
				drawPointer(display, pos, ili948x.WHITE)

				last_pos.Move(pos)
				redraw = false
			}

			// all done, wait for a while
			time.Sleep(100 * time.Millisecond)
		}
	}
}
