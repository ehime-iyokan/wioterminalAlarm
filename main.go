package main

import (
	"log"
	"time"
)

var (
	// wifiアクセスポイントの ssid, パスワードを設定する
	ssid     string
	password string
)

func main() {
	_, err := Wifi(ssid, password, 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)

	}
}
