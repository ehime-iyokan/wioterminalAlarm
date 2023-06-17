package main

import (
	"machine"
	"runtime"
	"time"

	"tinygo.org/x/drivers/net/http"
	"tinygo.org/x/drivers/rtl8720dn"
)

var (
	rtl       *rtl8720dn.Driver
	connected bool
	uart      UARTx
	debug     bool
	buf       [0x1000]byte
)

// SetupRTL8720DN sets up the RTL8270DN for use.
func SetupRTL8720DN() (*rtl8720dn.Driver, error) {
	rtl = rtl8720dn.New(machine.UART3, machine.PB24, machine.PC24, machine.RTL8720D_CHIP_PU)

	if debug {
		waitSerial()
	}

	rtl.Debug(debug)
	rtl.Configure()

	connected = true
	return rtl, nil
}

// Wifi sets up the RTL8720DN and connects it to Wi-Fi.
func AdjustTimeUsingWifi(ssid, pass string, timeout time.Duration) (*rtl8720dn.Driver, error) {
	_, err := SetupRTL8720DN()
	if err != nil {
		return nil, err
	}

	err = rtl.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil {
		return rtl, err
	}

	http.UseDriver(rtl)
	http.SetBuf(buf[:])

	// NTP
	t, err := GetCurrentTime()
	if err != nil {
		return nil, err
	}
	runtime.AdjustTimeOffset(-1 * int64(time.Since(t)))

	return rtl, nil
}

// Wait for user to open serial console
func waitSerial() {
	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}
}

type UARTx struct {
	*machine.UART
}

func fetchTimeNowJst() time.Time {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Now().UTC().In(jst)
	// nsは0で固定する。ns単位までは比較は行わないため
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

type Alarm struct {
	time         time.Time
	ringing      bool
	selectorTime int // alarm.selectorTime = 0:秒調整, 1:時間調整
}

// デフォルトのアラーム時刻を設定する。second, ns は 0 が設定される
func (a *Alarm) setDefaultTime(t time.Time) {
	a.time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

// selectorTime が 0 なら "分"、1 なら "秒"を +1 する
func (a *Alarm) timeIncrement() {
	if a.selectorTime == 0 {
		minuteIncrementer, _ := time.ParseDuration("1m")
		a.time = a.time.Add(minuteIncrementer)
	} else {
		hourIncrementer, _ := time.ParseDuration("1h")
		a.time = a.time.Add(hourIncrementer)
	}
}

// selectorTime が 0 なら "分"、1 なら "秒"を -1 する
func (a *Alarm) timeDecrement() {
	if a.selectorTime == 0 {
		minuteDecrementer, _ := time.ParseDuration("-1m")
		a.time = a.time.Add(minuteDecrementer)
	} else {
		hourDecrementer, _ := time.ParseDuration("-1h")
		a.time = a.time.Add(hourDecrementer)
	}
}

// 今の時間とアラームの時間が一致しているかを確認し、一致している場合は funcAlarmOn を実行する
func (a *Alarm) alarmOnIfTimeMatched(t_now time.Time, funcAlarmOn func()) {
	t_now = time.Date(t_now.Year(), t_now.Month(), t_now.Day(), t_now.Hour(), t_now.Minute(), t_now.Second(), 0, t_now.Location())
	if t_now.Equal(a.time) {
		a.ringing = true
		funcAlarmOn()
	}
}

// アラームをOFFにし、funcAlarmOff を実行する
func (a *Alarm) alarmOff(funcAlarmOff func()) {
	a.ringing = false
	funcAlarmOff()
}

// 日付が変わっても時間が一致するようにするため、年月日を今の時間に同期させる。
func (a *Alarm) adjustDay(t_now time.Time) {
	a.time = time.Date(t_now.Year(), t_now.Month(), t_now.Day(), a.time.Hour(), a.time.Minute(), a.time.Second(), 0, t_now.Location())
}
