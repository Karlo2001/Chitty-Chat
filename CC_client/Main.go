package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "CHITTY-CHAT/CC_proto"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	id         int32
)

// Registers the client with the server.
func Join(client pb.ChittyChatClient, req *pb.ParticipantInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_id, err := client.Join(ctx, req)
	if err != nil {
		log.Fatalf(err.Error())
	}
	id = _id.Id
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
	for {
		fmt.Printf("Enter name: ")
		sc := bufio.NewScanner(os.Stdin)
		if sc.Scan() {
			name = sc.Text()
		}
		if name != "*** Server" && name != "" {
			break
		}
		fmt.Println("*** Invalid name")
	}

	Join(client, &pb.ParticipantInfo{Name: name})

	stream, err := client.Recieve(context.Background(), &pb.ParticipantId{Id: id})
	if err != nil {
		log.Fatalf(err.Error())
	}

	connected := true

	// Goroutine recieving messages from server
	go func() {
		for connected {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf(msg.Name + ": " + msg.Msg)
		}
	}()

	// Goroutine reading input
	go func() {
		println("Enter //Leave to disconnect")
		var in string
		sc := bufio.NewScanner(os.Stdin)
		for {
			if sc.Scan() {
				in = sc.Text()
			}
			if in == "//Leave" {
				_, err = client.Leave(context.Background(), &pb.ParticipantId{Id: id})
				if err != nil {
					log.Fatalf(err.Error())
				} else {
					connected = false
					return
				}
			} else {
				go func() {
					_, err := client.Publish(context.Background(), &pb.Msg{Id: id, Msg: in})
					if err != nil {
						log.Fatalf(err.Error())
					}
				}()
			}
		}
	}()

	// Keep client running
	for connected {
		time.Sleep(1 * time.Second)
	}
}
