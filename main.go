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
	glay := color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	black := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}

	alarm := Alarm{}

	timeNow := time.Time{}
	timeNowBefore := time.Time{}

	// ハードウェア設定処理開始 ---------------------------------------------------------
	display := initdisplay.InitDisplay()
	display.FillScreen(white)

	_, err := AdjustTimeUsingWifi(ssid, password, 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	display.FillScreen(black)
	labelTimeNow := NewLabel(72, 320)
	SettingAlarmlabel := NewLabel(48, 320)

	pwm := machine.TCC0
	pwm.Configure(machine.PWMConfig{})
	channelA, _ := pwm.Channel(machine.BUZZER_CTR)
	pwm.SetPeriod(uint64(1e9) / 440)

	button_3 := machine.BUTTON_3
	button_2 := machine.BUTTON_2

	button_3.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_2.Configure(machine.PinConfig{Mode: machine.PinInput})

	crosskey := CrossKey{
		push:  machine.SWITCH_U,
		up:    machine.SWITCH_X,
		down:  machine.SWITCH_B,
		right: machine.SWITCH_Z,
		left:  machine.SWITCH_Y,
	}

	crosskey.up.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.down.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.left.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.right.Configure(machine.PinConfig{Mode: machine.PinInput})

	// ハードウェア設定処理終了 ---------------------------------------------------------

	// 割り込み処理設定開始 ------------------------------------------------------------
	button_3.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.mode ^= 1
	})
	button_2.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.ringing = false
		pwm.Set(channelA, 0)
	})

	crosskey.up.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.timeIncrement()
	})
	crosskey.down.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.timeDecrement()
	})
	crosskey.left.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.selectorTime ^= 1
	})
	crosskey.right.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.selectorTime ^= 1
	})
	// 割り込み処理設定終了 ------------------------------------------------------------

	// ループ処理開始 ------------------------------------------------------------------
	alarm.setDefaultTime(fetchTimeNowJst())

	for {
		timeNow = fetchTimeNowJst()

		if alarm.mode == 0 {
			// 時間表示モード
			if alarm.modeBefore == 1 {
				// 画面遷移直後の処理
				labelTimeNow.FillScreen(glay)
			}

			// ns単位までは比較は行わない
			timeNow = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second(), 0, timeNow.Location())
			if timeNow.Equal(alarm.time) {
				alarm.ringing = true
				pwm.Set(channelA, pwm.Top()/4)
			}

			stringTimeNow := fmt.Sprintf("%04d/%02d/%02d\n%02d:%02d:%02d",
				timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
			if alarm.ringing == true {
				stringTimeNow = stringTimeNow + "\n!Alarm-ON!"
			}

			if timeNow.Equal(timeNowBefore) {
				// 何もしない
			} else {
				// 情報に変化があれば表示内容を更新する
				labelTimeNow.FillScreen(glay)
				tinyfont.WriteLine(labelTimeNow, &freemono.Regular12pt7b, 0, 18, stringTimeNow, white)
				display.DrawRGBBitmap(0, 0, labelTimeNow.Buf, labelTimeNow.W, labelTimeNow.H)
			}
		} else {
			// 時間設定モード
			stringTimeAlarm := fmt.Sprintf("setting alarm\n%02d:%02d", alarm.time.Hour(), alarm.time.Minute())

			if alarm.modeBefore == 0 {
				// 画面遷移直後の処理
				display.FillScreen(black)
				SettingAlarmlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlarmlabel, &freemono.Regular12pt7b, 0, 18, stringTimeAlarm, white)
				display.DrawRGBBitmap(0, 0, SettingAlarmlabel.Buf, SettingAlarmlabel.W, SettingAlarmlabel.H)
			}

			if alarm.time.Equal(alarm.timeBefore) {
				// 何もしない
			} else {
				// 情報に変化があれば表示内容を更新する
				SettingAlarmlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlarmlabel, &freemono.Regular12pt7b, 0, 18, stringTimeAlarm, white)
				display.DrawRGBBitmap(0, 0, SettingAlarmlabel.Buf, SettingAlarmlabel.W, SettingAlarmlabel.H)
			}

			// 年月日を同期させる。日付が変わってもアラームが鳴るようにするため
			alarm.ajustDay(timeNow)
		}

		timeNowBefore = timeNow
		alarm.modeBefore = alarm.mode
		alarm.timeBefore = alarm.time

		time.Sleep(10 * time.Millisecond)
	}
}
