package grpc_captcha_validation

import "errors"

var (
	ERR_CAPTCHA_SECRET_IS_EMPTY               = errors.New("secret: secret is empty")
	ERR_FAILED_CREATE_REQUEST                 = errors.New("secret: failed to create http request")
	ERR_FAILED_FIND_META_DATA_WITH_HEADER_KEY = errors.New("captcha: failed to get metadata from context")
	ERR_FAILED_FIND_HEADER_IN_META_DATA       = errors.New("captcha: failed to find header key in metadata")
)
