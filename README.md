# GRPC captcha validator
The GRPC Captcha Validation Middleware is a Go module that provides a reusable middleware component for validating captcha challenges on GRPC servers. The middleware intercepts incoming requests, verifies the captcha challenge for the request using a configurable provider, and rejects the request if the challenge is not valid.

The module is designed to be configurable, with options to specify the captcha provider, challenge validation endpoint, and custom headers for the validation request. It uses the Google Protobuf library to parse the GRPC method descriptor and extract custom options defined using protobuf extensions.

## Feature
- support Google, Cloudflare, hcaptcha provider to verify challenge
- support unary and stream server middleware
- captcha validation for specific rpc methods

# Example

example captcha for GRPC server:

```protobuf
syntax = "proto3";
package example;
option go_package = "github.com/GoFarsi/grpc-captcha-validation/_example/proto";

import "captcha.proto";

service GreetingService {
  rpc Greeting(GreetingRequest) returns(GreetingResponse) {
    option(captcha.captcha) = {
      check_challenge: true,
      provider: GOOGLE,
    };
  }

  rpc GreetingStream(stream GreetingRequest) returns(stream GreetingResponse) {
    option(captcha.captcha) = {
      check_challenge: true,
      provider: GOOGLE,
    };
  }
}

message GreetingRequest {
  string name = 1;
}

message GreetingResponse {
  string message = 1;
}
```

```go
package main

import (
	"context"
	"fmt"
	gcv "github.com/GoFarsi/grpc-captcha-validation"
	"github.com/GoFarsi/grpc-captcha-validation/_example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
)

type GreetingService struct {
	proto.UnimplementedGreetingServiceServer
}

func (g *GreetingService) Greeting(ctx context.Context, in *proto.GreetingRequest) (*proto.GreetingResponse, error) {
	return &proto.GreetingResponse{
		Message: fmt.Sprintf("greeting %s", in.Name),
	}, nil
}

func (g *GreetingService) GreetingStream(srv proto.GreetingService_GreetingStreamServer) error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		resp := proto.GreetingResponse{Message: fmt.Sprintf("Greeting %s", req.Name)}
		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	c := gcv.NewCaptcha("foo", "", "", "")

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(c.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(c.StreamServerInterceptor()),
	)

	reflection.Register(srv)

	proto.RegisterGreetingServiceServer(srv, &GreetingService{})

	log.Fatalln(srv.Serve(listener))
}
```
