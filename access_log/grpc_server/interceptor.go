package grpc_server

import (
	"context"
	"time"

	"github.com/Mr-LvGJ/jota/log"
	"google.golang.org/grpc"

	"github.com/Mr-LvGJ/jota/access_log"
)

func AccessLogInterceptor(logger *access_log.AccessLogger, opts ...access_log.Option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var attrs []any
		attrs = append(attrs, "method", info.FullMethod)
		logger.Info(ctx, "Handle request begin")
		startTime := time.Now()
		resp, err := handler(ctx, req)
		cost := time.Since(startTime)
		if err != nil {
			log.Error(ctx, "Handle request err", "request", req, "cost", cost, "err", err)
			return resp, err
		}
		logger.Info(ctx, "Handle request done", "resp", resp, "cost", cost)
		return resp, err
	}
}
