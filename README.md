# environment
* windows10
* tinygo version 0.27.0 windows/amd64 (using go version go1.19 and LLVM version 15.0.0)
* wioterminal
# alerm
* 機能一覧
  * 現在時刻を表示する
    * 実装済
  * 時刻を設定する
    * ボタン3(※)を押すと設定画面に遷移
    * 実装済
  * 設定した時刻にブザーを鳴らす
    * ボタン2を押すとブザー停止
      * 実装済

 ※:ボタンはwioterminal上側のボタンであり、左から「ボタン3」,「ボタン2」,「ボタン1」
# require
* 以下のようなファイルをalermのディレクトリ内に格納してください
  * wifiアクセスポイントのSSIDとパスワードを「xxxx」,「yyyy」にそれぞれ入力する必要があります
```go
[ssid_init.go]
package main

func init() {
	ssid = "xxxx"
	password = "yyyy"
}
````
