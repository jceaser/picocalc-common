//go:build tinygo

/* **********************************************************************************************100
A very simple, uneducated interface for the sound system on the PicoCalc. This assumes two speakers,
Left and Right and that they are on pins 26 and 27. 

This code is directly inspired by:
https://github.com/laingcc/Picocalc-Coyote-OS/blob/master/pwm_sound/pwm_sound.c

Open Source AI (through ollama) was use to translate, but frankly it was so horrible I refuse to
give credit because of the work I had to do after and I don't want to say anything bad about them.
Ollama is great though.

Other sources of inspiration:
* https://tinygo.org/docs/tutorials/pwm/
************************************************************************************************* */

package sound

import (
	"machine"
	"sync"
	"time"
)

type SoundType int

const (
	SND_BEEP      SoundType = iota
	SND_TAB_SWITCH SoundType = iota
	SND_ERROR     SoundType = iota
)

type Sound struct {
	lock_left sync.Mutex
	lock_right sync.Mutex
	inuse_left bool
	inuse_right bool
}

const (
	AUDIO_PIN_L = machine.GP26 	// Replace with the actual GP pin for left audio
	AUDIO_PIN_R = machine.GP27 	// Replace with the actual GP pin for right audio
	PWM_FREQUENCY = 1000      	// Default frequency for sound in Hz
	PWM_PERIOD = 65_535			// Set your desired PWM clock frequency in kHz
)

// ***************************************************************************80
// MARK - Setup

func Initialize() Sound {
	gpioL := AUDIO_PIN_L
	gpioR := AUDIO_PIN_R

	// Configure the pins for PWM
	gpioL.Configure(machine.PinConfig{Mode: machine.PinPWM})
	gpioR.Configure(machine.PinConfig{Mode: machine.PinPWM})

	machine.PWM5.Configure(machine.PWMConfig{Period: PWM_PERIOD})
	sound := Sound{}
	return sound
}

// ***************************************************************************80
// MARK - unsafe functions

func validFrequncy(raw int) bool {
	return 20 <= raw && raw <= 20_000
}

func invalidFrequncy(raw int) bool {
	return raw < 20 || 20_000 < raw
}

/* play a pre defined sound */
func Play(snd SoundType) {
	var frequency int
	var duration time.Duration

	switch snd {
	case SND_BEEP:
		frequency = 1000
		duration = time.Millisecond * 100 // 100ms
	case SND_TAB_SWITCH:
		frequency = 1500
		duration = time.Millisecond * 50 // 50ms
	case SND_ERROR:
		frequency = 400
		duration = time.Millisecond * 200 // 200ms
	}

	SetPwmLevel(frequency, duration)
}

/* direct PWM control on both channels */
func SetPwmLevel(frequency int, duration time.Duration) {
	period := time.Second / time.Duration(frequency)
	dutyCycle := uint32(PWM_PERIOD / 2)

	// Set PWM period in nanoseconds
	machine.PWM5.Configure(machine.PWMConfig{Period: uint64(period.Nanoseconds())})

	left, _ := machine.PWM5.Channel(AUDIO_PIN_L)
	right, _ := machine.PWM5.Channel(AUDIO_PIN_R)

	machine.PWM5.Set(left, dutyCycle) // Set left channel
	machine.PWM5.Set(right, dutyCycle) // Set right channel
	time.Sleep(duration)
	machine.PWM5.Set(left, 0) // Stop left channel
	machine.PWM5.Set(right, 0) // Stop right channel
}

/* direct PWM control on the left channel */
func SetPwmLevelLeft(frequency int, duration time.Duration) {
	period := time.Second / time.Duration(frequency)
	dutyCycle := uint32(PWM_PERIOD / 2)

	// Set PWM period in nanoseconds
	machine.PWM5.Configure(machine.PWMConfig{Period: uint64(period.Nanoseconds())})

	left, _ := machine.PWM5.Channel(AUDIO_PIN_L)

	machine.PWM5.Set(left, dutyCycle) // Set left channel
	time.Sleep(duration)
	machine.PWM5.Set(left, 0) // Stop left channel
}

/* direct PWM control on the left channel */

func SetPwmLevelRight(frequency int, duration time.Duration) {
	period := time.Second / time.Duration(frequency)
	dutyCycle := uint32(PWM_PERIOD / 2)

	// Set PWM period in nanoseconds
	machine.PWM5.Configure(machine.PWMConfig{Period: uint64(period.Nanoseconds())})

	right, _ := machine.PWM5.Channel(AUDIO_PIN_R)

	machine.PWM5.Set(right, dutyCycle) // Set right channel
	time.Sleep(duration)
	machine.PWM5.Set(right, 0) // Stop right channel
}

// ***************************************************************************80
// MARK - safe functions

/*

These functions use a thread to ensure that a sound finishes before another
starts. Using these functions along with other functions may result in unexpected

*/

func (s *Sound) Play(snd SoundType) {
	s.lock_left.Lock()
	s.lock_right.Lock()
	if s.inuse_left || s.inuse_right {
		return
	}
	s.inuse_left = true
	s.inuse_right = true

	go func(which SoundType) {
		Play(which)
		s.inuse_left = false
		s.inuse_right = false
		s.lock_left.Unlock()
		s.lock_right.Unlock()
	}(snd)
}

func (s Sound) InUseRight() bool {
	return s.inuse_right
}

func (s Sound) InUseLeft() bool {
	return s.inuse_left
}

func (s Sound) PlayFrequency(frequency int, duration time.Duration) {
	if invalidFrequncy(frequency) {
		return
	}
	s.PlayFrequencyLeft(frequency, duration)
	s.PlayFrequencyRight(frequency, duration)
}

// Public method to play a sound as a thread
func (s *Sound) PlayFrequencyLeft(frequency int, duration time.Duration) {
	if invalidFrequncy(frequency) {
		return
	}
	if s.inuse_left {
		return
	}
	s.lock_left.Lock()
	s.inuse_left = true

	go func(freq int, dur time.Duration) {
		SetPwmLevelLeft(freq, dur)
		s.inuse_left = false
		s.lock_left.Unlock()
	}(frequency, duration)
}

// Public method to play a sound as a thread
func (s *Sound) PlayFrequencyRight(frequency int, duration time.Duration) {
	if invalidFrequncy(frequency) {
		return
	}
	if s.inuse_right {
		return
	}
	s.lock_right.Lock()
	s.inuse_right = true

	go func(freq int, dur time.Duration) {
		SetPwmLevelRight(freq, dur)
		s.inuse_right = false
		s.lock_right.Unlock()
	}(frequency, duration)
}
