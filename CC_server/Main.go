package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "CHITTY-CHAT/CC_proto"

	"google.golang.org/grpc"
)

var (
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of courses")
	port       = flag.Int("port", 10000, "The server port")
)

type chittyChatServer struct {
	pb.UnimplementedCoursesServer
}

type client struct {
	id   int32
	name string
}

func (s *chittyChatServer) Join(ctx context.Context, in *pb.CourseRequest) (*pb.Course, error) {

	// No course was found, return nil
	return nil, nil
}

func (s *chittyChatServer) Recieve(id *pb.ParticipantId, srv pb.StreamService) {

}

/*
func NewServer() *ChittyChatServer {
	s := &chittyChatServer{}
	return s
}
*/

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChittyChatServer(grpcServer, chittyChatServer)
	grpcServer.Serve(lis)
}
