Ecommerce Review Service

电商评论微服务
---
## 技术栈
参数校验 bufbuild/protovalidate 

gRPC 生成 REST api grpc-ecosystem/

grpc-gateway 

生成 CRUD 操作代码 go-gorm/gen 

## Use
### Docker
```bash
$ docker compose -f docker-compose.dev.yml up -d
```

### Make
导入 db.sql 到 MySQL
```bash
$ cp .env.example .env
$ nano .env
$ make dev
```
