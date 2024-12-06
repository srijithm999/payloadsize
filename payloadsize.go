package payloadsize

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func init() {
	caddy.RegisterModule(PayloadSize{})
}

type PayloadSize struct {
	MaxPayloadSize int `json:"max_payload_size,omitempty"` // Optional: Maximum allowed payload size in bytes
	logger         *zap.Logger
}

func (j PayloadSize) Provision(context caddy.Context) error {
	j.logger = context.Logger()
	return nil
}

// CaddyModule returns the Caddy module information
func (PayloadSize) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.payloadsize",
		New: func() caddy.Module { return new(PayloadSize) },
	}
}

func (j *PayloadSize) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	j.logger.Info("JWTPayload Serve Http processing request")
	// Read the body (plaintext payload)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return caddyhttp.Error(http.StatusBadRequest, fmt.Errorf("failed to read request body: %v", err))
	}

	payloadSize := len(body)
	fmt.Printf("Request body payload size: %d bytes\n", payloadSize)

	// Enforce max payload size if specified
	if j.MaxPayloadSize > 0 && payloadSize > j.MaxPayloadSize {
		return caddyhttp.Error(http.StatusRequestEntityTooLarge, fmt.Errorf("JWT payload size exceeds maximum allowed size of %d bytes", j.MaxPayloadSize))
	}

	j.logger.Info("Recorded Payload size", zap.Int("size", payloadSize))

	// Proceed to the next handler
	return next.ServeHTTP(w, r)
}

var (
	_ caddy.Module                = (*PayloadSize)(nil)
	_ caddyhttp.MiddlewareHandler = (*PayloadSize)(nil)
	_ caddy.Provisioner           = (*PayloadSize)(nil)
	//_ caddyfile.Unmarshaler       = (*JWTPayload)(nil)
)
