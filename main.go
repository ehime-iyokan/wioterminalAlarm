package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"github.com/sago35/tinydisplay/examples/initdisplay"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

var (
	// wifiアクセスポイントの ssid, パスワードを設定する
	ssid     string
	password string
)

type CrossKey struct {
	push  machine.Pin
	up    machine.Pin
	down  machine.Pin
	right machine.Pin
	left  machine.Pin
}

func main() {
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	display := initdisplay.InitDisplay()
	display.FillScreen(white)
	_, err := AdjustTime(ssid, password, 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	glay := color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	black := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	display.FillScreen(black)
	labelTimeNow := NewLabel(72, 320)
	SettingAlermlabel := NewLabel(48, 320)

	// mode = 0:時間設定モード, 1:時間表示モード
	mode := 1
	button_3 := machine.BUTTON_3
	button_3.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_3.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		mode ^= 1
	})

	pwm := machine.TCC0
	pwm.Configure(machine.PWMConfig{})
	channelA, _ := pwm.Channel(machine.BUZZER_CTR)
	pwm.SetPeriod(uint64(1e9) / 440)

	AlermON := 0
	button_2 := machine.BUTTON_2
	button_2.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_2.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		AlermON = 0
		pwm.Set(channelA, 0)
	})

	crosskey := CrossKey{
		push:  machine.SWITCH_U,
		up:    machine.SWITCH_X,
		down:  machine.SWITCH_B,
		right: machine.SWITCH_Z,
		left:  machine.SWITCH_Y,
	}
	// flgAlermSetting = 0:秒調整, 1:時間調整
	flgAlermSetting := 0
	alermMinute := 0
	alermHour := 0

	crosskey.up.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.up.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if flgAlermSetting == 0 {
			if 0 <= alermMinute && alermMinute < 59 {
				alermMinute++
			} else {
				alermMinute = 0
			}
		} else {
			if 0 <= alermHour && alermHour < 23 {
				alermHour++
			} else {
				alermHour = 0
			}
		}
	})
	crosskey.down.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.down.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if flgAlermSetting == 0 {
			if 0 < alermMinute && alermMinute <= 59 {
				alermMinute--
			} else {
				alermMinute = 59
			}
		} else {
			if 0 < alermHour && alermHour <= 23 {
				alermHour--
			} else {
				alermHour = 23
			}
		}

	})
	crosskey.left.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.left.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		flgAlermSetting ^= 1
	})
	crosskey.right.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.right.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		flgAlermSetting ^= 1
	})

	timeNow := fetchStringNowJst()
	timeNowBefore := timeNow
	timeAlermStringBefore := ""
	modeBefore := 0

	for {
		timeNow = fetchStringNowJst()

		if mode == 1 {
			if modeBefore == 0 {
				labelTimeNow.FillScreen(glay)
			}

			timeNowString := fmt.Sprintf("%04d/%02d/%02d\n%02d:%02d:%02d",
				timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())

			if timeNow.Hour() == alermHour && timeNow.Minute() == alermMinute {
				AlermON = 1
				alermHour = 88
				alermMinute = 88
				pwm.Set(channelA, pwm.Top()/4)
			}

			if AlermON == 1 {
				timeNowString = timeNowString + "\n!Alerm-ON!"
			}

			if timeNow.Second() == timeNowBefore.Second() {
				// 何もしない
			} else {
				labelTimeNow.FillScreen(glay)
				tinyfont.WriteLine(labelTimeNow, &freemono.Regular12pt7b, 0, 18, timeNowString, white)
				display.DrawRGBBitmap(0, 0, labelTimeNow.Buf, labelTimeNow.W, labelTimeNow.H)
			}
		} else {
			timeAlermString := fmt.Sprintf("setting alerm\n%02d:%02d", alermHour, alermMinute)

			if modeBefore == 1 {
				display.FillScreen(black)
				SettingAlermlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlermlabel, &freemono.Regular12pt7b, 0, 18, timeAlermString, white)
				display.DrawRGBBitmap(0, 0, SettingAlermlabel.Buf, SettingAlermlabel.W, SettingAlermlabel.H)
			}

			if timeAlermString == timeAlermStringBefore {
				// 何もしない
			} else {
				SettingAlermlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlermlabel, &freemono.Regular12pt7b, 0, 18, timeAlermString, white)
				display.DrawRGBBitmap(0, 0, SettingAlermlabel.Buf, SettingAlermlabel.W, SettingAlermlabel.H)
			}

			timeAlermStringBefore = timeAlermString
		}

		modeBefore = mode
		timeNowBefore = timeNow

		time.Sleep(10 * time.Millisecond)
	}
}
