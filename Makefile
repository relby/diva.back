genproto:
	protoc --go_out=./pkg/genproto --go_opt=paths=source_relative -I ./api/proto --go-grpc_out=./pkg/genproto --go-grpc_opt=paths=source_relative ./api/proto/customers.proto

gensqlc:
	sqlc -f ./database/sqlc.yaml generate