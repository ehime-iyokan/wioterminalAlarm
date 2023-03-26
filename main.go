package main

import (
	"fmt"
	"log"
	"time"
)

var (
	// wifiアクセスポイントの ssid, パスワードを設定する
	ssid     string
	password string
)

func main() {
	_, err := AdjustTime(ssid, password, 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	for {
		time.Sleep(1 * time.Second)
		now := time.Now()
		nowUTC := now.UTC()
		nowJST := nowUTC.In(jst)

		fmt.Println(nowJST)
	}
}
