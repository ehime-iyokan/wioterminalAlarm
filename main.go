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

	for {
		time.Sleep(1 * time.Second)
		timeNow := fetchNowJst()
		fmt.Println(timeNow)
	}
}

func fetchNowJst() time.Time {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now()
	nowUTC := now.UTC()
	return nowUTC.In(jst)
}
