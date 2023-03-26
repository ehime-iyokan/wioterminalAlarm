# wioterminalのwifiモジュールを扱う

## wifiモジュールの設定
* rtl8720dn.go より
  * input: wifiアクセスポイントの ssid, pass, timeoutとなる時間
  * output:エラーのみ(time.Now()で現在時刻が取得できるようになる)

## 参考
* "github.com/sago35/tinygo-examples/wioterminal/initialize"
  * rtl8720dn.go
  * ntp.go
