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
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
