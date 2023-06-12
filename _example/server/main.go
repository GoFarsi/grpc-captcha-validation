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
