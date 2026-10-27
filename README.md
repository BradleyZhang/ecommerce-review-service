TODO
- [ ] 设计错误码(iam那个)


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