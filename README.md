# AnimeAPI Server (GAE/golang)

## 概要

Go言語製のAnimeAPIサーバーです。 [scala-playframework版](https://github.com/Project-ShangriLa/sora-playframework-scala)と互換性があります。

Google App Engine(GAE)で動作するように作られています。

マスターデータは従来のMySQLではなくGoogle CloudDatastoreで管理するので既存のAPIサーバーからデータマイグレーションする必要があります。[ツール](https://github.com/Project-ShangriLa/anime_master_migrate_google_datastore)


## セットアップ

```
cp app.sample.yaml app.yaml
```

app.yamlのapplication:を変更してください

## 実行

```
goapp serve
```

## デプロイ

```
goapp deploy
```

## ライブラリインストール

```
go get XXXX
```


## コードフォーマット

```
gofmt anime_api.go
gofmt -w anime_api.go
```

## lint

```
golint anime_api.go
```

## Go Pathサンプル

参考情報

```
export GOPATH=/Users/XXXX/gopath
export PATH=$PATH:~/software/go_appengine/:~/gopath/bin/
```