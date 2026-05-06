# picocalc-common
Common code for using go with the Clockwork PicoCalc

## Overview

This project was created as a way to formalize the code found at [rpngo][rpngo] and [tinygo-picocalc][tinygo-picocalc]. The goal of this libarry is to provide a basic level of interaction with the PicoCalc hardware (screen, keyboard, battery).

The code seems to be related to the code at :

## Code

There is a minimalistic main() function which does not do much. The package is inteded to be a library to pull into your projects and has the following imports:

* `github.com/jceaser/picocalc-common/i2ckbd` - everything through the i2c, like the keyboard
* `github.com/jceaser/picocalc-common/ili948x` - the screen
* `github.com/jceaser/picocalc-common/util` - general code

Example usage:

	//go:build tinygo
	package main
	import (
		"github.com/jceaser/picocalc-common/util"
		"github.com/jceaser/picocalc-common/ili948x"
	)
	func main() {
		if util.Version() >= 0.1 {
			ili948x.InitDisplay()
		}
	}

## Other projects

* https://github.com/clockworkpi/PicoCalc/tree/master/Code/picocalc_keyboard
* https://github.com/tinygo-org/tinyfont/blob/release/tinyfont.go
* https://pkg.go.dev/github.com/tinygo-org/tinygo/src/machine

## Credit

Most of this code I got from [MatttWach@github][matt-watch]. His example project [tinygo-picocalc][tinygo-picocalc] was the primary reason I bought a PicoCalc. I wanted to be able to do what he did and port some of my less polished [tools][gotools-rpn] over and see if they worked.

The battery code was inspired by [Cayote-OS][c-battry].

----

License: None as of yet. Many of these sources I read do not have a license, so I'm inclined to go with something pretty open but friendly to anyone who wants to try to earn a living. So the MIT License is what I plan to apply here as I have used that in some of my own code.

[matt-watch]: https://github.com/mattwach/
[rpngo]: https://github.com/mattwach/rpngo
[tinygo-picocalc]: https://github.com/mattwach/tinygo_picocalc
[gotools-rpn]: https://github.com/jceaser/gotools/blob/master/rpn.go
[c-battry]: https://github.com/laingcc/Picocalc-Coyote-OS/blob/master/i2ckbd/i2ckbd.c
