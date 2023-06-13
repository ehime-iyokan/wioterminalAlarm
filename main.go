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

	// mode = 0:時間設定モード, 1:時間表示モード
	mode := 1
	modeBefore := 0
	alarmRinging := 0

	crosskey := CrossKey{
		push:  machine.SWITCH_U,
		up:    machine.SWITCH_X,
		down:  machine.SWITCH_B,
		right: machine.SWITCH_Z,
		left:  machine.SWITCH_Y,
	}

	// selectorAlarmTime = 0:秒調整, 1:時間調整
	selectorAlarmTime := 0
	timeAlarm := time.Time{}
	timeAlarmBefore := time.Time{}
	minuteIncrementer, _ := time.ParseDuration("1m")
	hourIncrementer, _ := time.ParseDuration("1h")
	minuteDecrementer, _ := time.ParseDuration("-1m")
	hourDecrementer, _ := time.ParseDuration("-1h")

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

	button_3 := machine.BUTTON_3
	button_3.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_3.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		mode ^= 1
	})

	pwm := machine.TCC0
	pwm.Configure(machine.PWMConfig{})
	channelA, _ := pwm.Channel(machine.BUZZER_CTR)
	pwm.SetPeriod(uint64(1e9) / 440)

	button_2 := machine.BUTTON_2
	button_2.Configure(machine.PinConfig{Mode: machine.PinInput})
	button_2.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		alarmRinging = 0
		pwm.Set(channelA, 0)
	})

	crosskey.up.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.up.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if selectorAlarmTime == 0 {
			timeAlarm = timeAlarm.Add(minuteIncrementer)
		} else {
			timeAlarm = timeAlarm.Add(hourIncrementer)
		}
	})
	crosskey.down.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.down.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if selectorAlarmTime == 0 {
			timeAlarm = timeAlarm.Add(minuteDecrementer)
		} else {
			timeAlarm = timeAlarm.Add(hourDecrementer)
		}

	})
	crosskey.left.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.left.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		selectorAlarmTime ^= 1
	})
	crosskey.right.Configure(machine.PinConfig{Mode: machine.PinInput})
	crosskey.right.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		selectorAlarmTime ^= 1
	})

	// ハードウェア設定処理終了 ---------------------------------------------------------

	// ループ処理開始 ------------------------------------------------------------------
	timeAlarm = fetchTimeDefaultAlarmTime()

	for {
		timeNow = fetchTimeNowJst()

		if mode == 1 {
			// 時間表示モード
			if modeBefore == 0 {
				// 画面遷移直後の処理
				labelTimeNow.FillScreen(glay)
			}

			// ns単位までは比較は行わない
			timeNow = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second(), 0, time.FixedZone("Asia/Tokyo", 9*60*60))
			if timeNow.Equal(timeAlarm) {
				alarmRinging = 1
				pwm.Set(channelA, pwm.Top()/4)
			}

			timeNowString := fmt.Sprintf("%04d/%02d/%02d\n%02d:%02d:%02d",
				timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
			if alarmRinging == 1 {
				timeNowString = timeNowString + "\n!Alarm-ON!"
			}

			if timeNow.Equal(timeNowBefore) {
				// 何もしない
			} else {
				// 情報に変化があれば表示内容を更新する
				labelTimeNow.FillScreen(glay)
				tinyfont.WriteLine(labelTimeNow, &freemono.Regular12pt7b, 0, 18, timeNowString, white)
				display.DrawRGBBitmap(0, 0, labelTimeNow.Buf, labelTimeNow.W, labelTimeNow.H)
			}
		} else {
			// 時間設定モード
			timeAlarmString := fmt.Sprintf("setting alarm\n%02d:%02d", timeAlarm.Hour(), timeAlarm.Minute())

			if modeBefore == 1 {
				// 画面遷移直後の処理
				display.FillScreen(black)
				SettingAlarmlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlarmlabel, &freemono.Regular12pt7b, 0, 18, timeAlarmString, white)
				display.DrawRGBBitmap(0, 0, SettingAlarmlabel.Buf, SettingAlarmlabel.W, SettingAlarmlabel.H)
			}

			if timeAlarm.Equal(timeAlarmBefore) {
				// 何もしない
			} else {
				// 情報に変化があれば表示内容を更新する
				SettingAlarmlabel.FillScreen(glay)
				tinyfont.WriteLine(SettingAlarmlabel, &freemono.Regular12pt7b, 0, 18, timeAlarmString, white)
				display.DrawRGBBitmap(0, 0, SettingAlarmlabel.Buf, SettingAlarmlabel.W, SettingAlarmlabel.H)
			}

			// 年月日は同期させる。日付が変わってもアラームが鳴るようにするため
			timeAlarm = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), timeAlarm.Hour(), timeAlarm.Minute(), timeAlarm.Second(), 0, time.FixedZone("Asia/Tokyo", 9*60*60))
		}

		modeBefore = mode
		timeNowBefore = timeNow
		timeAlarmBefore = timeAlarm

		time.Sleep(10 * time.Millisecond)
	}
}
