TODO
- [x] 自动生成openapi.yml
- [ ] proto 抽出到独立 git 库，服务通过 gitsubmodule 调用
- [ ] 实现 b 端的 Reply
- [ ] 服务注册服务发现，负载均衡带discovery 的服务发现（grpc.Dial("consul//....")
- [ ] canal 部署 https://liwenzhou.com/posts/go/canal
- [ ] fafka 部署
- [ ] canal 配置投递到kafka

---
依赖
docker run --name redis507 -p 6379:6379 -d redis:5.0.7
docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root1234 -d mysql


---
biz(定义repo interface) 核心业务逻辑
model 核心业务实体

data(实现repo)->biz,model
service->biz,model

---
bufbuild/protovalidate 做参数校验
grpc-ecosystem/grpc-gateway gRPC 生成 REST api
go-gorm/gen 生成 CRUD 操作代码
bwmarrin/snowflake 做ID生成

---
trade off
虽然biz里嵌入了gerr 这个依赖gRPC 的错误函数，但若之后需删除gRPC，只需修改错误函数就是，没必要在抽象出错误interface


---
为什么使用 es：因为管理后台需要各种 模糊查询，在数据量大后，速度慢。以及查看的往往是联表的数据，在db里需要联表去查。这种场景往往用es


数据变更，增删改-》用数据库，查询用es。 db数据用 canal 和 kafka 发到es。用 CQRS 架构
前提是不需要即时的一致性，评论往往需要审核，对即时没有要求

singleflight 防止缓存击穿