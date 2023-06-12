package grpc_captcha_validation

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func extractHeaderValueFromContext(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ERR_FAILED_FIND_META_DATA_WITH_HEADER_KEY
	}

	foundedHeaders, ok := md[header]
	if !ok {
		return "", ERR_FAILED_FIND_HEADER_IN_META_DATA
	}

	return foundedHeaders[0], nil
}
