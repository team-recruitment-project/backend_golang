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

## アプリケーション
**起動**
```shell
go run ./cmd/teamrecruitment/main.go  
```

**全体テスト**
```shell
go test ./... -v
```

## API 明細書
OpenAPI形式で提供
```text
/api/*.yaml
```

### Auth
- [Auth API Specification](/api/auth.yaml)

### Team

- [Teams API Specification](/api/teams.yaml)
