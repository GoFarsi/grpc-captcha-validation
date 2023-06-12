package grpc_captcha_validation

import (
	"context"
	"fmt"
	"strings"
)

func challenge(ctx context.Context, providerAddress string, secretKey secret, customHeaderKey string) error {
	headerKey := _default_captcha_key

	if len(customHeaderKey) != 0 {
		headerKey = customHeaderKey
	}

	challengeKey, err := extractHeaderValueFromContext(ctx, headerKey)
	if err != nil {
		return err
	}

	resp := new(CaptchaResponse)
	parameter := map[string]string{
		"secret":   string(secretKey),
		"response": challengeKey,
	}

	if err := clientRequest(providerAddress, "POST", parameter, nil, resp); err != nil {
		return err
	}

	if !resp.Success {
		errorCode := strings.Join(resp.ErrorCodes, " ")
		return fmt.Errorf("captcha [%s]: your challenge is unsuccessful, try again to complete captcha challenge", errorCode)
	}

	return nil
}
