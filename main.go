package main

import (
	"image/color"
	"log"
	"machine"
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
	YMDlabel.FillScreen(glay)
	HMSlabel.FillScreen(glay)

	// mode = 0:時間設定モード, 1:時間表示モード
	mode := 1
	button_3 := machine.BUTTON_3
	button_3.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_3.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		mode ^= 1
	})

	hms_before := ""
	for {
		if mode == 1 {
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
			// ここに時間を設定する処理を書く
			display.FillScreen(black)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
