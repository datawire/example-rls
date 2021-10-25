package main

import (
	"context"
	"os"

	"github.com/datawire/dlib/dhttp"
	"github.com/datawire/dlib/dlog"
	"google.golang.org/grpc"

	envoyCoreV3 "github.com/datawire/ambassador/v2/pkg/api/envoy/config/core/v3"
	envoyRateLimitV3 "github.com/datawire/ambassador/v2/pkg/api/envoy/service/ratelimit/v3"
)

func main() {
	ctx := context.Background()

	grpcHandler := grpc.NewServer()
	envoyRateLimitV3.RegisterRateLimitServiceServer(grpcHandler, &RateLimitService{})

	sc := &dhttp.ServerConfig{
		Handler: grpcHandler,
	}

	dlog.Info(ctx, "starting...")
	if err := sc.ListenAndServe(ctx, ":3000"); err != nil {
		dlog.Errorf(ctx, "server exited with error: %v", err)
		os.Exit(1)
	}
	dlog.Info(ctx, "server exited without error")
}

type RateLimitService struct{}

func (s *RateLimitService) ShouldRateLimit(ctx context.Context, req *envoyRateLimitV3.RateLimitRequest) (*envoyRateLimitV3.RateLimitResponse, error) {
	return &envoyRateLimitV3.RateLimitResponse{
		OverallCode: envoyRateLimitV3.RateLimitResponse_OVER_LIMIT,
		RawBody:     []byte("Oh no, you have been rate-limted :(\n"),
		ResponseHeadersToAdd: []*envoyCoreV3.HeaderValue{
			{Key: "content-type", Value: "text/plain"},
		},
	}, nil
}
