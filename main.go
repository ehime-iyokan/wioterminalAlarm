package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"strconv"
	"strings"
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
	YMDlabel := NewLabel(24, 320)
	HMSlabel := NewLabel(24, 320)
	SettingAlermlabel := NewLabel(48, 320)

	// mode = 0:時間設定モード, 1:時間表示モード
	mode := 1
	button_3 := machine.BUTTON_3
	button_3.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_3.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		mode ^= 1
	})

	crosskey := CrossKey{
		push:  machine.SWITCH_U,
		up:    machine.SWITCH_X,
		down:  machine.SWITCH_B,
		right: machine.SWITCH_Z,
		left:  machine.SWITCH_Y,
	}
	// alermtime_select_flg = 0:秒調整, 1:時間調整
	alermtime_select_flg := 0
	alermtime_second := 0
	alermtime_hour := 0

	crosskey.up.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.up.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if alermtime_select_flg == 0 {
			if alermtime_second < 59 {
				alermtime_second++
			} else {
				alermtime_second = 0
			}
		} else {
			if alermtime_hour < 23 {
				alermtime_hour++
			} else {
				alermtime_hour = 0
			}
		}
	})
	crosskey.down.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.down.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if alermtime_select_flg == 0 {
			if alermtime_second > 0 {
				alermtime_second--
			} else {
				alermtime_second = 59
			}
		} else {
			if alermtime_hour > 0 {
				alermtime_hour--
			} else {
				alermtime_hour = 23
			}
		}

	})
	crosskey.left.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.left.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alermtime_select_flg ^= 1
	})
	crosskey.right.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.right.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alermtime_select_flg ^= 1
	})

	mode_before := 0
	hms_before := ""
	alermtime_string_before := ""

	for {
		if mode == 1 {
			if mode_before == 0 {
				YMDlabel.FillScreen(glay)
				HMSlabel.FillScreen(glay)
			}

			timeNow := fetchStringNowJst()
			array := strings.Split(timeNow, " ")
			ymd := array[0]
			hms := array[1]

			if hms_before == hms {
				// 何もしない
			} else {
				YMDlabel.FillScreen(glay)
				tinyfont.WriteLine(YMDlabel, &freemono.Regular12pt7b, 0, 18, ymd, white)
				display.DrawRGBBitmap(0, 0, YMDlabel.Buf, YMDlabel.W, YMDlabel.H)
				HMSlabel.FillScreen(glay)
				tinyfont.WriteLine(HMSlabel, &freemono.Regular12pt7b, 0, 18, hms, white)
				display.DrawRGBBitmap(0, 24, HMSlabel.Buf, HMSlabel.W, HMSlabel.H)
			}

			hms_before = hms

		} else {
			str_hour := strconv.Itoa(alermtime_hour)
			str_second := strconv.Itoa(alermtime_second)
			alermtime_string := fmt.Sprintf("setting alerm\n%02s:%02s", str_hour, str_second)

			if mode_before == 1 {
				display.FillScreen(black)
				SettingAlermlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlermlabel, &freemono.Regular12pt7b, 0, 18, alermtime_string, white)
				display.DrawRGBBitmap(0, 0, SettingAlermlabel.Buf, SettingAlermlabel.W, SettingAlermlabel.H)
			}

			if alermtime_string == alermtime_string_before {
				// 何もしない
			} else {
				SettingAlermlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlermlabel, &freemono.Regular12pt7b, 0, 18, alermtime_string, white)
				display.DrawRGBBitmap(0, 0, SettingAlermlabel.Buf, SettingAlermlabel.W, SettingAlermlabel.H)
			}

			alermtime_string_before = alermtime_string
		}

		mode_before = mode

		time.Sleep(10 * time.Millisecond)
	}
}
