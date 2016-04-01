package biolviewity

import (
	"fmt"
	"log"
	"net"

	pb "github.com/kitschysynq/biolviewity/setserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Set interface {
	Add(string) (bool, error)
	Remove(string) (bool, error)
	Contains(string) (bool, error)
}

type UnimplSet struct{}

var ErrNotImpl = fmt.Errorf("not implemented")

func (u *UnimplSet) Add(i string) (bool, error) {
	return false, ErrNotImpl
}

func (u *UnimplSet) Remove(i string) (bool, error) {
	return false, ErrNotImpl
}

func (u *UnimplSet) Contains(i string) (bool, error) {
	return false, ErrNotImpl
}

type SetServer struct {
	addr string
	set  Set
}

func NewSetServer(addr string, s Set) *SetServer { return &SetServer{addr: addr, set: s} }

var DefaultSetServer = NewSetServer(":8080", &UnimplSet{})

func (s *SetServer) Add(ctx context.Context, i *pb.Item) (*pb.Response, error) {
	present, err := s.set.Add(i.Item)
	if err != nil {
		return &pb.Response{
			Type:    pb.Response_ERROR,
			Present: false,
			Error:   err.Error(),
		}, nil
	}
	return &pb.Response{
		Type:    pb.Response_SUCCESS,
		Present: present,
		Error:   "",
	}, nil
}

func (s *SetServer) Remove(context.Context, *pb.Item) (*pb.Response, error) {
	return &pb.Response{
		Type:    pb.Response_ERROR,
		Present: false,
		Error:   ErrNotImpl.Error(),
	}, nil
}

func (s *SetServer) Contains(context.Context, *pb.Item) (*pb.Response, error) {
	return &pb.Response{
		Type:    pb.Response_ERROR,
		Present: false,
		Error:   ErrNotImpl.Error(),
	}, nil
}

func (s *SetServer) AddHandler(h Set) {
	s.set = h
}

func (s *SetServer) ListenAndServe() {
	addr := s.addr
	if addr == "" {
		addr = ":8080"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g := grpc.NewServer()
	pb.RegisterSetServer(g, s)
	g.Serve(l)
}

func AddHandler(h Set) {
	DefaultSetServer.AddHandler(h)
}

func ListenAndServe() {
	DefaultSetServer.ListenAndServe()
}
