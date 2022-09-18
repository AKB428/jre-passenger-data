# JR東 乗車人員のデータ分析ツール

データ元 https://www.jreast.co.jp/passenger/

## 概要
JR東日本のWEB公開されている各駅の乗車人員データ2000-2021年度までのデータを取得し
CSV形式で出力するツールです

## コンパイル

```
go build -trimpath -ldflags '-s -w' -o ./bin/jre-passenger-data main.go
```

```
go build -trimpath -ldflags '-s -w' -o ./bin/download ./download/download.go
```

※Windows向けバイナリの場合は

```
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags '-s -w' -o ./bin/jre-passenger-data.exe main.go
```

```
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags '-s -w' -o ./bin/download.exe ./download/download.go
```


## 実行

負荷を考慮してダウンロードプロセスとパースプロセスを分離しています

### htmlデータのダウンロード

```
./bin/download
```
実行カレントディレクトリにhtmlsフォルダが作成されますので、そのデータを後続の処理で使用します

### 乗車人員データTOP100の出力

```
./bin/jre-passenger-data > csv/count.csv
```

### 乗車人員データTOP100(ランキング)の出力

```
./bin/jre-passenger-data -r > csv/rank.csv
```