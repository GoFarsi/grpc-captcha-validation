gen:
	protoc --proto_path=captcha --proto_path=captcha/google --go_out=. --go-grpc_out=. captcha/*.proto

gen-test:
	protoc --proto_path=_example/proto --proto_path=captcha --go_out=. --go-grpc_out=. _example/proto/*.proto