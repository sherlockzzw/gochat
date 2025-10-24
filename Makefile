GROUP_NAME=gochat
PROJECT_NAME=gochat
API_ROOT=/api/v1
API_PROTO_FILES=$(shell find api -name *.proto)
API_PB_FILES=$(shell find api -name "*.pb.go")

grpc:
	protoc --proto_path=./api \
	   	--proto_path=./third_party \
	   	--go_out=paths=source_relative:./api \
	   	--go-grpc_out=paths=source_relative:./api \
		--validate_out=paths=source_relative,lang=go:./api \
		--experimental_allow_proto3_optional \
		$(API_PROTO_FILES)
	$(MAKE) inject-tags

# 在pb.go文件中注入自定义tag
# go install github.com/favadi/protoc-go-inject-tag@latest
.PHONY: inject-tags
inject-tags:
	$(foreach file, $(API_PB_FILES), protoc-go-inject-tag -input=$(file);)

# 安装必要的工具
install-tools:
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

# 清理生成的文件
clean:
	find api -name "*.pb.go" -delete
	find api -name "*_grpc.pb.go" -delete
	find api -name "*.pb.validate.go" -delete
	find api -name "*.swagger.json" -delete

.PHONY: grpc inject-tags install-tools clean
