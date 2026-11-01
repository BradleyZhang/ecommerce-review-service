.PHONY: api gen build dev
%: build
gen:
	@go run ./cmd/gen
	@go generate ./...

build:gen
	@mkdir -p bin/ && go build -o ./bin/review-service ./cmd/review-service

dev:build
	@echo "===== build finish ====="
	@./bin/review-service

canal:
	./run-canal.sh -e canal.instance.master.address=host.docker.internal:3306 \
         -e canal.instance.dbUsername=canal \
         -e canal.instance.dbPassword=canal \
         -e canal.instance.connectionCharset=UTF-8 \
         -e canal.instance.tsdb.enable=true \
         -e canal.instance.gtidon=false \
		 -e canal.mq.topic=default-topic \
 		 -e canal.mq.dynamicTopic=.*\\..* \

# 之后需要 vi canal-server/conf/canal.properties,修改并重启
#  canal.serverMode = kafka
# kafka.bootstrap.servers = host.docker.internal:29092 # 按具体修改


consul:
	docker run -d \
		--name consul \
		--restart unless-stopped \
		-p 8500:8500 \
		-p 8600:8600/udp \
		-v consul_data:/consul/data \
		hashicorp/consul:1.21.3 \
		agent \
		-server \
		-bootstrap-expect=1 \
		-node=consul-server \
		-datacenter=dc1 \
		-ui \
		-client=0.0.0.0 \
		-data-dir=/consul/data