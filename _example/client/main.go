package main

import (
	"context"
	"fmt"
	"github.com/GoFarsi/grpc-captcha-validation/_example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

var name = []string{
	"Javad",
	"Ali",
	"Reza",
	"Ahmad",
}

func main() {
	dialOpts := make([]grpc.DialOption, 0)
	dialOpts = append(dialOpts, grpc.WithInsecure())

	conn, err := grpc.Dial(":8080", dialOpts...)
	if err != nil {
		log.Fatalln(err)
	}

	g := proto.NewGreetingServiceClient(conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-captcha-key", "bar")

	resp, err := g.Greeting(ctx, &proto.GreetingRequest{
		Name: "Javad",
	})

	if err == nil {
		fmt.Println(resp.Message)
	}

	stream, err := g.GreetingStream(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer stream.CloseSend()

	for _, n := range name {
		req := &proto.GreetingRequest{Name: n}

		if err := stream.Send(req); err != nil {
			log.Fatalln(err)
		}
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(msg)
	}

}
