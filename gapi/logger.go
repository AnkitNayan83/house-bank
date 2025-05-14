package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {

	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown

	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status", statusCode.String()).
		Dur("duration", duration).
		Msg("recieved grpc request")

	return result, err
}

// override the default http.ResponseWriter to capture the status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	Body       []byte
}

func (rec *ResponseWriter) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *ResponseWriter) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseWriter{ResponseWriter: res, statusCode: http.StatusOK}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()

		if rec.statusCode >= 400 {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", "http").
			Int("status_code", rec.statusCode).
			Str("status", http.StatusText(rec.statusCode)).
			Str("method", req.Method).
			Str("path", req.URL.Path).
			Dur("duration", duration).
			Msg("recieved http request")
	})
}
