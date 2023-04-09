# environment
* windows10
* tinygo version 0.26.0 windows/amd64 (using go version go1.19.3 and LLVM version 14.0.0)
* wioterminal
# alerm
* 機能一覧
  * 現在時刻を表示する
  * 時刻を設定する
    * ボタン3(※)を押すと設定画面に遷移
  * 設定した時刻にブザーを鳴らす
    * ボタン2を押すとブザー停止

 ※:ボタンはwioterminal上側のボタンを指し、左から「ボタン3」,「ボタン2」,「ボタン1」
# require
* 以下のようなファイルをalermのディレクトリ内に格納してください
  * 「xxxx」等は適したものに書き換える
```go
[ssid_init.go]
package main

func init() {
	ssid = "xxxx"
	password = "yyyy"
}
````
