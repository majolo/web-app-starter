package services

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/majolo/web-app-starter/gen/diary/v1"
	"github.com/nedpals/supabase-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Entry struct {
	text      string
	createdAt time.Time
	userId    string
}

type Service struct {
	inMemDiary map[int64]Entry
	diary.UnimplementedDiaryServiceServer
	supabaseClient *supabase.Client
}

func NewDiaryService(supabase *supabase.Client) *Service {
	s := &Service{
		inMemDiary:     map[int64]Entry{0: {text: "starter entry", createdAt: time.Now(), userId: "1"}},
		supabaseClient: supabase,
	}
	return s
}

func (s *Service) CreateEntry(ctx context.Context, req *diary.CreateEntryRequest) (*diary.CreateEntryResponse, error) {
	user, err := s.verifyTokenSupabase(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetText() == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}
	entry := Entry{
		text:      req.GetText(),
		createdAt: time.Now(),
		userId:    user.ID,
	}
	id := int64(len(s.inMemDiary) + 1)
	s.inMemDiary[id] = entry
	return &diary.CreateEntryResponse{
		Id: id,
	}, nil
}

func (s *Service) ListEntries(ctx context.Context, req *diary.ListEntriesRequest) (*diary.ListEntriesResponse, error) {
	user, err := s.verifyTokenSupabase(ctx)
	if err != nil {
		return nil, err
	}
	var entries []*diary.Entry
	for id, entry := range s.inMemDiary {
		if entry.userId != user.ID {
			continue
		}
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

func (s *Service) verifyTokenSupabase(ctx context.Context) (*supabase.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata")
	}
	cookiesStr, ok := md["grpcgateway-cookie"]
	if !ok {
		return nil, fmt.Errorf("no cookie")
	}
	// New code to extract sAccessToken
	var sAccessToken string
	for _, cookieStr := range cookiesStr {
		header := http.Header{}
		header.Add("Cookie", cookieStr)
		request := http.Request{Header: header}
		cookies := request.Cookies()
		for _, cookie := range cookies {
			// find the cookie that starts with sb- and ends with -auth-token
			// this is not sustainable as these names are internal but it can be fixed longer term
			if strings.HasPrefix(cookie.Name, "sb-") &&
				strings.HasSuffix(cookie.Name, "-auth-token") {
				sAccessToken = cookie.Value
				break
			}
		}
	}
	if sAccessToken == "" {
		return nil, fmt.Errorf("no sAccessToken found")
	}
	decodedValue, err := url.QueryUnescape(sAccessToken)
	if err != nil {
		return nil, err
	}
	token := extractSupabaseToken(decodedValue)
	user, err := s.supabaseClient.Auth.User(ctx, token)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Note this is fairly hardcoded, Supabase documentation isn't great around this and doesn't seem to provide a nice way to handle the cookies they produce.
// This should be revisited.
func extractSupabaseToken(cookie string) string {
	splits := strings.Split(cookie, "\"")
	if len(splits) < 2 {
		return ""
	}
	return splits[1]
}
