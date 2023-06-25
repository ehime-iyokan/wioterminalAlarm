# environment
* windows11
* tinygo version 0.27.0 windows/amd64 (using go version go1.19.9 and LLVM version 15.0.0)
* wioterminal
# function
* 機能一覧
  * 時刻表示機能
    * 現在時刻を表示する
      * 実装済
  * アラーム機能
    * 時刻を設定する
      * ボタン3(※)を押して設定画面に遷移する
        * 実装済
      * 十字キーで 時間/分 を設定する
        * 上下キーで時刻選択する
          * 実装済
        * 左右キーで 時間/分 を切り替える
          * 実装済
    * 設定した時刻にブザーを鳴らす
      * 実装済
    * ボタン2を押してブザーを停止する
      * 実装済

 ※:ボタンはwioterminal上側のボタンであり、左から「ボタン3」,「ボタン2」,「ボタン1」
# require
* 以下のようなファイルをalarmのディレクトリ内に格納してください
  * wifiアクセスポイントのSSIDとパスワードを「xxxx」,「yyyy」にそれぞれ入力する必要があります
```go
[ssid_init.go]
package main

func init() {
	ssid = "xxxx"
	password = "yyyy"
}
````
