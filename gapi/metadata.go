package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mdata := &Metadata{}

	// 1. fetch the userAgent and IP from the grpc metadata package when request is coming through GRPC Gateway
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			mdata.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			mdata.UserAgent = userAgent[0]
		}

		if clientIP := md.Get(xForwardedForHeader); len(clientIP) > 0 {
			mdata.ClientIP = clientIP[0]
		}
	}

	// 2. fetch the client IP from peer package when the request is through peer package
	if pr, ok := peer.FromContext(ctx); ok {
		mdata.ClientIP = pr.Addr.String()
	}

	return mdata
}
