package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/sago35/tinydisplay/examples/initdisplay"
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

	for {
		time.Sleep(1 * time.Second)
		timeNow := fetchStringNowJst()
		fmt.Println(timeNow)
	}
}

func fetchStringNowJst() string {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now()
	nowUTC := now.UTC()
	nowJST := nowUTC.In(jst)
	return nowJST.Format("2006/01/02 15:04:05")
}
