package httphelper

import (
	"companies-service/config"
	"companies-service/pkg/logger"
	"companies-service/pkg/sanitize"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetRequestID(ctx *gin.Context) string {
	// Get the request ID is in the headers
	requestID := ctx.Request.Header.Get("X-Request-ID")

	return requestID
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// Get ctx with timeout and request id from echo context
func GetCtxWithReqID(c *gin.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetRequestID(c))
	return ctx, cancel
}

// Get context  with request id
func GetRequestCtx(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// Configure jwt cookie
func ConfigureJWTCookie(cfg *config.Config, jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:       cfg.Cookie.Name,
		Value:      jwtToken,
		Path:       "/",
		RawExpires: "",
		MaxAge:     cfg.Cookie.MaxAge,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}

// ReadRequest gets body and validate
func ReadRequest(c *gin.Context, request interface{}) error {
	if err := c.Bind(request); err != nil {
		return err
	}

	validate := validator.New()

	return validate.StructCtx(c, request)
}

// SanitizeRequest sanitizes and validates request
func SanitizeRequest(c *gin.Context, request interface{}) error {
	if err := c.Bind(&request); err != nil {
		return err
	}

	sanBody, err := sanitize.SanitizeJSON(request)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(sanBody, request); err != nil {
		return err
	}

	validate := validator.New()

	return validate.StructCtx(c, request)
}

// ErrResponseWithLog writes error response with logging error for gin context
func ErrResponseWithLog(c *gin.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(c),
		c.ClientIP(),
		err,
	)

	c.JSON(ErrorResponse(err))
}

// LogResponseError logs the error response with logging error for gin context
func LogResponseError(ctx *gin.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		ctx.ClientIP(),
		err,
	)
}
