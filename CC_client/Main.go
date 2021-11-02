package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	pb "CHITTY-CHAT/CC_proto"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	id         int32
)

type lClock struct {
	t  int32
	mu sync.Mutex
}

var clock = lClock{}

func updateClock(time int32) {
	clock.mu.Lock()
	if clock.t < time {
		clock.t = time
	}
	clock.t++
	clock.mu.Unlock()
}
func incrementClock() {
	clock.mu.Lock()
	clock.t++
	clock.mu.Unlock()
}

// Registers the client with the server.
func Join(client pb.ChittyChatClient, req *pb.ParticipantInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_id, err := client.Join(ctx, req)
	if err != nil {
		log.Fatalf(err.Error())
	}
	id = _id.Id
	updateClock(_id.Time)
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

	Join(client, &pb.ParticipantInfo{Time: clock.t, Name: name})

	incrementClock()
	stream, err := client.Stream(context.Background(), &pb.ParticipantId{Time: clock.t, Id: id})
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
			updateClock(msg.Time)
			if msg.Name == "*** Server" {
				log.Println(msg.Msg, clock.t)
			} else {
				log.Println("(" + strconv.Itoa(int(clock.t)) + ", " + msg.Name + "): " + msg.Msg)
			}
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
				if in == "//Leave" {
					incrementClock()
					_, err = client.Leave(context.Background(), &pb.ParticipantId{Time: clock.t, Id: id})
					if err != nil {
						log.Fatalf(err.Error())
					} else {
						connected = false
						return
					}
				} else {
					go func() {
						incrementClock()
						_, err := client.Publish(context.Background(), &pb.MsgFromClient{Time: clock.t, Id: id, Msg: in})
						if err != nil {
							log.Fatalf(err.Error())
						}
					}()
				}
			}

		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_, err = client.Leave(context.Background(), &pb.ParticipantId{Time: clock.t, Id: id})
		if err != nil {
			log.Fatalf(err.Error())
		}
		os.Exit(1)
	}()

	// Keep client running
	for connected {
		time.Sleep(1 * time.Second)
	}
}
