package main

import (
	"context"
	pb "github.com/anton.okolelov/otus-calendar-grpc/internal/grpc"
	"github.com/anton.okolelov/otus-calendar-grpc/internal/model"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"net"
)

//go:generate protoc --go_out=plugins=grpc:. internal/grpc/calendar.proto

type Server struct {
	pb.UnimplementedCalendarServer
	calendar *model.Calendar
}

func NewServer(calendar *model.Calendar) *Server {
	return &Server{
		calendar: calendar,
	}
}

func (server *Server) AddEvent(ctx context.Context, req *pb.Event) (*pb.EventId, error) {
	start, err := ptypes.Timestamp(req.Start)
	if err != nil {
		return nil, err
	}

	end, err := ptypes.Timestamp(req.End)
	if err != nil {
		return nil, err
	}
	id := uint32(server.calendar.AddEvent(model.Event{start, end, req.Payload}))
	return &pb.EventId{
		Id: id,
	}, nil
}

func (server *Server) UpdateEvent(ctx context.Context, req *pb.EventUpdateInfo) (*pb.EventId, error) {
	start, err := ptypes.Timestamp(req.Event.Start)
	if err != nil {
		return nil, err
	}

	end, err := ptypes.Timestamp(req.Event.End)
	if err != nil {
		return nil, err
	}
	server.calendar.UpdateEvent(model.EventId(req.EventId), model.Event{start, end, req.Event.Payload})
	return &pb.EventId{
		Id: req.EventId,
	}, nil
}

func (server *Server) DeleteEvent(ctx context.Context, req *pb.EventId) (*pb.EventId, error) {
	server.calendar.DeleteEvent(model.EventId(req.Id))
	return req, nil
}

func (server *Server) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.EventList, error) {
	events := server.calendar.Events()

	list := pb.EventList{
		Events: []*pb.Event{},
	}

	for _, event := range events {
		start, err := ptypes.TimestampProto(event.Start)
		if err != nil {
			return nil, err
		}
		end, err := ptypes.TimestampProto(event.End)

		if err != nil {
			return nil, err
		}

		list.Events = append(list.Events, &pb.Event{
			Start: start,
			End:   end,
		})
	}
	return &list, nil
}

const (
	port = ":50051"
)

func main() {
	log.Printf("Run server at port %v", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalendarServer(s, NewServer(model.New()))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
