package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "CHITTY-CHAT/CC_proto"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	id         int32
)

func Join(client pb.chittyChatClient, req *pb.ParticipantInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_id, err := client.Join(ctx, req)
	if err != nil {
		// Error handling
		return
	}
	id = _id
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewChittyChatClient(conn)

	var name string
	fmt.Printf("Enter name: ")
	fmt.Scanln(&name)

	Join(client, &pb.ParticipantInfo{time: 1, name: name})
	// Get stream
	stream, err := client.Recieve(&pb.ParticipantId{id})
	if err != nil {
		// Error handling
	}

	// Goroutine recieving messages from server
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// Stream ended
				return
			}
			if err != nil {
				//Error handling
			}
			log.Printf(msg)
		}
	}()

}
