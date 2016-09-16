package crdt_server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/criprooff/crdt_server/setserver"
	serf "github.com/hashicorp/serf/serf"
	"github.com/satori/go.uuid"
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

type Member struct {
	Name   string
	Addr   string
	Status string
}

type SetServer struct {
	addr string
	set  Set

	mu   sync.Mutex
	serf *serf.Serf
}

const (
	setServerPrefix = "setserver"
)

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

func (s *SetServer) ListenAndServe() error {
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
	return g.Serve(l)
}

func (s *SetServer) RegisterSelf(addrs []string, bindAddr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.serf != nil {
		return errors.New("already registered")
	}

	host, _, err := net.SplitHostPort(bindAddr)
	if err != nil {
		return err
	}

	cfg := serf.DefaultConfig()
	cfg.NodeName = strings.Join([]string{setServerPrefix, uuid.NewV4().String()}, "-")
	cfg.MemberlistConfig.BindAddr = host
	cfg.MemberlistConfig.AdvertiseAddr = host

	s.serf, err = serf.Create(cfg)
	if err != nil {
		return err
	}

	_, err = s.serf.Join(addrs, true)
	return err
}

func (s *SetServer) DeregisterSelf() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	defer s.serf.Shutdown()

	if err := s.serf.Leave(); err != nil {
		return err
	}

	s.serf = nil
	return nil
}

func (s *SetServer) Peers() []Member {
	members := make([]Member, 0, len(s.serf.Members()))

	for _, m := range s.serf.Members() {
		if !strings.HasPrefix(m.Name, setServerPrefix) {
			continue
		}

		portStr := strconv.Itoa(int(m.Port))

		members = append(members, Member{
			Name:   m.Name,
			Addr:   net.JoinHostPort(m.Addr.String(), portStr),
			Status: m.Status.String(),
		})
	}

	return members
}

func AddHandler(h Set) {
	DefaultSetServer.AddHandler(h)
}

func ListenAndServe() error {
	go func() {
		for range time.Tick(10 * time.Second) {
			for _, m := range DefaultSetServer.Peers() {
				log.Println("member:", m)
			}
		}
	}()
	return DefaultSetServer.ListenAndServe()
}

func RegisterSelf(addrs []string, bindAddr string) error {
	return DefaultSetServer.RegisterSelf(addrs, bindAddr)
}

func DeregisterSelf() error {
	return DefaultSetServer.DeregisterSelf()
}
