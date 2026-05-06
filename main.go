//go:build tinygo

/************************************************************************************************100
A RPN calculator based on the gotools calculator I built some time ago

Compile with:
	 tinygo build -o tinygo_test.uf2 -target=pico
	 tinygo build -o tinygo_test2.uf2 -target=pico2
***************************************************************************************************/

package main

import (
	/*"bytes"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"*/

	//"picocalc/i2ckbd"
	//"picocalc/ili948x"

	//"tinygo.org/x/tinyfont"
	//"tinygo.org/x/tinyfont/freemono"

	"github.com/jceaser/picocalc-common/lib"
	"github.com/jceaser/picocalc-common/ili948x"
	"machine"
)

// 320x320 screen size
func main() {
	if lib.Version() == 0.1 {
		machine.Watchdog.Start()
		ili948x.InitDisplay()
	}
}