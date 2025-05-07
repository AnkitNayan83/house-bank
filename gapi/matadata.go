package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// Metadata keys
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	xForwardedFor        = "x-forwarded-for"
	userAgent            = "user-agent"
)

type Metadata struct {
	ClientIp  string
	UserAgent string
}

func (server *Server) extractMetaData(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val := md[grpcGatewayUserAgent]; len(val) > 0 {
			mtdt.UserAgent = val[0]
		} else if val := md[userAgent]; len(val) > 0 {
			mtdt.UserAgent = val[0]
		} else {
			log.Println("User-Agent not found in metadata")
		}

		if val := md[xForwardedFor]; len(val) > 0 {
			mtdt.ClientIp = val[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		if p.Addr != nil {
			mtdt.ClientIp = p.Addr.String()
		}
	}

	return mtdt
}
