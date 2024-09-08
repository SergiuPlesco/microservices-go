
protomsg:
	protoc --go_out=. --go_opt=paths=import \
    --go-grpc_out=. --go-grpc_opt=paths=import \
	proto/currency.proto

    


proto_grpc: 
	protoc -I proto/ proto/currency.proto --go_out_grpc=proto/currency

proto_msg:
	protoc --go_out=internal/adapters/framework/left/grpc --proto_path=internal/adapters/framework/left/grpc/proto/ \
    internal/adapters/framework/left/grpc/proto/*_msg.proto --experimental_allow_proto3_optional && \
    protoc-go-inject-tag -input=./internal/adapters/framework/left/grpc/pb/*_msg.pb.go

proto_svc:
	protoc --go-grpc_out=require_unimplemented_servers=false:internal/adapters/framework/left/grpc \
	--proto_path=internal/adapters/framework/left/grpc/proto \
    internal/adapters/framework/left/grpc/proto/records_services.proto


		# protoc --go_out=./pb --go_opt=paths=import \
    # --go-grpc_out=./pb --go-grpc_opt=paths=import \
    # proto/currency.proto