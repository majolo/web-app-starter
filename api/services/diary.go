package services

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/majolo/web-app-starter/gen/diary/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Entry struct {
	text      string
	createdAt time.Time
}

type Service struct {
	inMemDiary map[int64]Entry
	diary.UnimplementedDiaryServiceServer
}

func NewDiaryService() *Service {
	s := &Service{
		inMemDiary: map[int64]Entry{},
	}
	return s
}

func (s *Service) CreateEntry(ctx context.Context, req *diary.CreateEntryRequest) (*diary.CreateEntryResponse, error) {
	entry := Entry{
		text:      req.GetText(),
		createdAt: time.Now(),
	}
	id := int64(len(s.inMemDiary) + 1)
	s.inMemDiary[id] = entry
	return &diary.CreateEntryResponse{
		Id: id,
	}, nil
}

func (s *Service) ListEntries(ctx context.Context, req *diary.ListEntriesRequest) (*diary.ListEntriesResponse, error) {
	var entries []*diary.Entry
	for id, entry := range s.inMemDiary {
		entries = append(entries, &diary.Entry{
			Id:        id,
			Text:      entry.text,
			CreatedAt: timestamppb.New(entry.createdAt),
		})
	}
	return &diary.ListEntriesResponse{
		Entries: entries,
	}, nil
}

func (s *Service) RegisterGRPC(gs *grpc.Server) {
	diary.RegisterDiaryServiceServer(gs, s)
}

func (s *Service) RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux) {
	err := diary.RegisterDiaryServiceHandlerServer(ctx, mux, s)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
}
