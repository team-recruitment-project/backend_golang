# Team Recruitment Project
サイドプロジェクトや個人開発ではなくチーム員と一緒に何かを作ったり勉強会を開くためのチーム員を募集するためのサイトです。

## 技術スタック
- Go
- Gin
- Ent
- MySQL
- Docker

## Ent 
**Schema 生成**
```shell
go run -mod=mod entgo.io/ent/cmd/ent new [SchemaName]
```

**Query 関連ファイル生成**
```shell
go generate ./ent
```