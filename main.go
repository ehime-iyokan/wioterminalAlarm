package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	"github.com/ehime-iyokan/alarm"
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

	alarm := alarm.Alarm{}
	mode := 0 // mode = 0:時間表示モード, 1:時間設定モード

	// ハードウェア設定処理開始 ---------------------------------------------------------
	display := initdisplay.InitDisplay()
	display.FillScreen(white)

	_, err := AdjustTimeUsingWifi(ssid, password, 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	display.FillScreen(black)
	labelTime := NewLabel(72, 320)

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
		mode ^= 1
	})
	button_2.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.AlarmOff(func() { pwm.Set(channelA, 0) })
	})

	crosskey.up.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.TimeIncrement()
	})
	crosskey.down.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarm.TimeDecrement()
	})
	crosskey.left.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		valueToggled := alarm.GetStatusSelectorTime() ^ 1
		alarm.SetSelectorTime(valueToggled)
	})
	crosskey.right.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		valueToggled := alarm.GetStatusSelectorTime() ^ 1
		alarm.SetSelectorTime(valueToggled)
	})
	// 割り込み処理設定終了 ------------------------------------------------------------

	// メインの処理開始 ------------------------------------------------------------------
	alarm.SetDefaultTime(fetchTimeNowJst())

	modeBefore := 0
	timeAlarmBefore := time.Time{}
	timeNowBefore := time.Time{}

	for {
		timeNow := fetchTimeNowJst()

		if mode == 0 {
			// 時間表示モード
			if modeBefore == 1 {
				// 画面遷移直後の処理
				labelTime.FillScreen(glay)
			}

			stringTimeNow := fmt.Sprintf("%04d/%02d/%02d\n%02d:%02d:%02d",
				timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())

			alarm.AlarmOnIfTimeMatched(timeNow, func() { pwm.Set(channelA, pwm.Top()/4) })

			if alarm.GetStatusRinging() == true {
				stringTimeNow = stringTimeNow + "\n!Alarm-ON!"
			}

			if timeNow.Equal(timeNowBefore) {
				// 何もしない
			} else {
				// 情報に変化があれば表示内容を更新する
				labelTime.FillScreen(glay)
				tinyfont.WriteLine(labelTime, &freemono.Regular12pt7b, 0, 18, stringTimeNow, white)
				display.DrawRGBBitmap(0, 0, labelTime.Buf, labelTime.W, labelTime.H)
			}
		} else {
			// 時間設定モード
			stringTimeAlarm := fmt.Sprintf("setting alarm\n%02d:%02d", alarm.GetTime().Hour(), alarm.GetTime().Minute())

			if modeBefore == 0 {
				// 画面遷移直後の処理
				labelTime.FillScreen(glay)
				tinyfont.WriteLine(labelTime, &freemono.Regular12pt7b, 0, 18, stringTimeAlarm, white)
				display.DrawRGBBitmap(0, 0, labelTime.Buf, labelTime.W, labelTime.H)
			}

			if alarm.GetTime().Equal(timeAlarmBefore) {
				// 何もしない
			} else {
				// 情報に変化があれば表示内容を更新する
				labelTime.FillScreen(glay)
				tinyfont.WriteLine(labelTime, &freemono.Regular12pt7b, 0, 18, stringTimeAlarm, white)
				display.DrawRGBBitmap(0, 0, labelTime.Buf, labelTime.W, labelTime.H)
			}

			alarm.AdjustDay(timeNow)
		}

		timeNowBefore = timeNow
		modeBefore = mode
		timeAlarmBefore = alarm.GetTime()

		time.Sleep(10 * time.Millisecond)
	}

}
