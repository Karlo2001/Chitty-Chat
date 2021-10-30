package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "CHITTY-CHAT/CC_proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type chittyChatServer struct {
	pb.UnimplementedChittyChatServer
}

type client struct {
	id   int32
	name string
	left bool
	ch   chan message
}

var clients []client

type message struct {
	clientId int32
	name     string
	text     string
}

// Searches for an inactive client in the "clients" slice, to be replaced.
// If none is found, the new client is added to the end of the slice.
// Also sends a join message to each active client.
func (s *chittyChatServer) Join(ctx context.Context, in *pb.ParticipantInfo) (*pb.ParticipantId, error) {
	var id int32
	var added bool
	for i := 0; i < len(clients); i++ {
		if clients[i].left {
			id = int32(i)
			clients[i] = client{id: id, name: in.Name, ch: make(chan message, 100)}
			added = true
			break
		}
	}
	if !added {
		id = int32(len(clients))
		nclient := client{id: id, name: in.Name, ch: make(chan message, 100)}
		clients = append(clients, nclient)
	}

	log.Println("Joined id: ", id, ", name: ", in.Name)
	for _, client := range clients {
		if !client.left {
			client.ch <- message{name: "*** Server", text: in.Name + " has joined!"}
		}
	}
	return &pb.ParticipantId{Id: id}, nil
}

// Updates the "Left" status on the civen client.
// The clients is then registered as no longer active.
// Also sends a leave message to each active client.
func (s *chittyChatServer) Leave(ctx context.Context, id *pb.ParticipantId) (*emptypb.Empty, error) {
	clients[id.Id].left = true

	log.Println("Left id: ", id.Id, ", name: ", clients[id.Id].name)
	for _, client := range clients {
		if !client.left {
			client.ch <- message{name: "*** Server", text: clients[id.Id].name + " has left the chat."}
		}
	}
	return &emptypb.Empty{}, nil
}

// Waits for a message to enter the channel belonging to the client.
// When a message is received, it is sent to the client through the server.
func (s *chittyChatServer) Recieve(id *pb.ParticipantId, srv pb.ChittyChat_RecieveServer) error {
	for !clients[id.Id].left {
		var msg = <-clients[id.Id].ch
		err := srv.Send(&pb.MsgRecieved{Name: msg.name, Msg: msg.text})
		if err != nil {
			log.Print(err.Error())
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// Sends the published message to the channels
// belonging to each active client, except the sender.
func (s *chittyChatServer) Publish(ctx context.Context, msg *pb.Msg) (*emptypb.Empty, error) {
	m := message{clientId: msg.Id, name: clients[msg.Id].name, text: msg.Msg}
	for _, client := range clients {
		if !client.left && msg.Id != client.id {
			client.ch <- m
		}
	}
	return &emptypb.Empty{}, nil
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
	grpcServer.Serve(lis)
}
