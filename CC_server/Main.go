package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "CHITTY-CHAT/CC_proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = flag.Int("port", 10000, "The server port")
	id   int32
)

type chittyChatServer struct {
	pb.UnimplementedChittyChatServer
}

type client struct {
	id   int32
	name string
	left bool
	ch   chan *pb.Msg
}

var clients []client

type lClock struct {
	t  []int32
	mu sync.Mutex
}

var clock = lClock{}

func incrementClock() {
	clock.mu.Lock()
	clock.t[id]++
	clock.mu.Unlock()
}

// Searches for an inactive client in the "clients" slice, to be replaced.
// If none is found, the new client is added to the end of the slice.
// Also calls Broadcast to send a join message to each active client.
func (s *chittyChatServer) Join(ctx context.Context, in *pb.ParticipantInfo) (*pb.ParticipantId, error) {
	var id int32
	var added bool
	for i := 1; i < len(clients); i++ {
		if clients[i].left {
			id = int32(i)
			clients[i] = client{id: id, name: in.Name, ch: make(chan *pb.Msg, 100)}
			added = true
			break
		}
	}
	if !added {
		id = int32(len(clients))
		nclient := client{id: id, name: in.Name, ch: make(chan *pb.Msg, 100)}
		clients = append(clients, nclient)
	}
	updateClock(in.Time)

	m := "*** Participant " + in.Name + " joined Chitty-Chat at Lamport time" // + strconv.Itoa(int(clock.t))
	Broadcast(&pb.Msg{Name: "*** Server", Msg: m})
	return &pb.ParticipantId{Time: clock.t, Id: id}, nil
}

// Updates the "Left" status on the civen client.
// The client is then registered as no longer active.
// Also calls Broadcast to send a leave message to each active client.
func (s *chittyChatServer) Leave(ctx context.Context, id *pb.ParticipantId) (*emptypb.Empty, error) {
	updateClock(id.Time)
	clients[id.Id].left = true

	m := "*** Participant " + clients[id.Id].name + " left Chitty-Chat at Lamport time" // + strconv.Itoa(int(clock.t))
	Broadcast(&pb.Msg{Name: "*** Server", Msg: m})
	return &emptypb.Empty{}, nil
}

// Waits for a message to enter the channel belonging to the client.
// When a message is received, it is sent to the client through the stream.
func (s *chittyChatServer) Stream(id *pb.ParticipantId, srv pb.ChittyChat_StreamServer) error {
	for !clients[id.Id].left {
		var msg = <-clients[id.Id].ch
		err := srv.Send(msg)
		if err != nil {
			log.Print(err.Error())
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// Receives the published message and sends it to the Broadcast method.
func (s *chittyChatServer) Publish(ctx context.Context, msg *pb.Msg) (*emptypb.Empty, error) {
	updateClock(msg.Time)
	m := pb.Msg{Name: clients[msg.Id].name, Msg: msg.Msg}
	Broadcast(&m)
	return &emptypb.Empty{}, nil
}

// Writes message to the server,
// before sending it to the channels belonging to each client.
func Broadcast(msg *pb.Msg) {
	incrementClock()
	msg.Time = clock.t
	if msg.Name == "*** Server" {
		log.Println(msg.Msg, clock.t)
	} else {
		log.Println(clock.t, msg.Name+": "+msg.Msg)
	}
	for _, client := range clients {
		if !client.left && client.id != 0 {
			client.ch <- msg
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChittyChatServer(grpcServer, &chittyChatServer{})
	id = 0
	clock.t = make([]int32, id+1)
	nclient := client{id: 0, name: "", left: false}
	clients = append(clients, nclient)
	grpcServer.Serve(lis)
}

func max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func lamportTime(clock []int32) int32 {
	var l int32
	for _, s := range clock {
		l += s
	}
	return l
}

func updateClock(otherClock []int32) {
	clock.mu.Lock()
	for i, s := range otherClock {
		if i >= len(clock.t) {
			clock.t = append(clock.t, s)
		} else {
			clock.t[i] = max(clock.t[i], s)
		}
	}
	clock.mu.Unlock()
	incrementClock()
}
