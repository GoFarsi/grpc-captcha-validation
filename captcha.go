package grpc_captcha_validation

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoFarsi/grpc-captcha-validation/captcha"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

const (
	_default_captcha_key = "x-captcha-key"
)

const (
	_default_google_validate_address     = "https://www.google.com/recaptcha/api/siteverify"
	_default_cloudflare_validate_address = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	_default_hcaptcha_validate_address   = "https://hcaptcha.com/siteverify"
)

type (
	secret string
)

type Captcha struct {
	googleSecret     secret
	cloudflareSecret secret
	hCaptchaSecret   secret
	customHeader     string
}

type CaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs any      `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Credit      bool     `json:"credit"`
	Score       float32  `json:"score"`
	ScoreReason []any    `json:"score_reason"`
}

// NewCaptcha create captcha object for access to middleware method
func NewCaptcha(googleRecaptchaSecret, cloudflareRecaptchaSecret, hCaptchaSecret, customHeaderKey string) *Captcha {
	return &Captcha{
		googleSecret:     secret(googleRecaptchaSecret),
		cloudflareSecret: secret(cloudflareRecaptchaSecret),
		hCaptchaSecret:   secret(hCaptchaSecret),
		customHeader:     customHeaderKey,
	}
}

// UnaryServerInterceptor captcha validator for unary server interceptor
func (c *Captcha) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		cp, err := c.extractCaptchaOptionFromDescriptor(strings.Replace(info.FullMethod[1:], "/", ".", -1))
		if err != nil {
			return nil, err
		}

		if cp.CheckChallenge {
			switch cp.Provider {
			case captcha.Provider_GOOGLE:
				if err = c.googleSecret.validate(); err != nil {
					return nil, err
				}
				if err = challenge(ctx, _default_google_validate_address, c.googleSecret, c.customHeader); err != nil {
					return nil, err
				}
			case captcha.Provider_CLOUDFLARE:
				if err = c.cloudflareSecret.validate(); err != nil {
					return nil, err
				}
				if err = challenge(ctx, _default_cloudflare_validate_address, c.cloudflareSecret, c.customHeader); err != nil {
					return nil, err
				}
			case captcha.Provider_HCAPTCHA:
				if err = c.hCaptchaSecret.validate(); err != nil {
					return nil, err
				}
				if err = challenge(ctx, _default_hcaptcha_validate_address, c.hCaptchaSecret, c.customHeader); err != nil {
					return nil, err
				}
			}
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor captcha validator for stream server interceptor
func (c *Captcha) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		cp, err := c.extractCaptchaOptionFromDescriptor(strings.Replace(info.FullMethod[1:], "/", ".", -1))
		if err != nil {
			return err
		}

		if cp.CheckChallenge {
			switch cp.Provider {
			case captcha.Provider_GOOGLE:
				if err = c.googleSecret.validate(); err != nil {
					return err
				}
				if err = challenge(stream.Context(), _default_google_validate_address, c.googleSecret, c.customHeader); err != nil {
					return err
				}
			case captcha.Provider_CLOUDFLARE:
				if err = c.cloudflareSecret.validate(); err != nil {
					return err
				}
				if err = challenge(stream.Context(), _default_cloudflare_validate_address, c.cloudflareSecret, c.customHeader); err != nil {
					return err
				}
			case captcha.Provider_HCAPTCHA:
				if err = c.hCaptchaSecret.validate(); err != nil {
					return err
				}
				if err = challenge(stream.Context(), _default_hcaptcha_validate_address, c.hCaptchaSecret, c.customHeader); err != nil {
					return err
				}
			}
		}

		return handler(srv, stream)
	}
}

func (c *Captcha) extractCaptchaOptionFromDescriptor(methodName string) (*captcha.Captcha, error) {
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(methodName))
	if err != nil {
		return nil, err
	}

	methodDesc, ok := desc.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, errors.New("captcha: failed to assertion method descriptor")
	}

	options, ok := methodDesc.Options().(*descriptorpb.MethodOptions)
	if !ok {
		return nil, errors.New("captcha: failed to assertion descriptor options")
	}

	capExt, err := proto.GetExtension(options, captcha.E_Captcha)
	if err != nil {
		fmt.Errorf("captcha: failed to get proto extension, got error %s", err.Error())
	}

	cp, ok := capExt.(*captcha.Captcha)
	if !ok {
		return nil, errors.New("captcha: failed to assertion captcha object with proto extension")
	}

	return cp, nil
}

func (c secret) validate() error {
	if len(c) == 0 {
		return ERR_CAPTCHA_SECRET_IS_EMPTY
	}
	return nil
}
