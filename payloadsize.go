package payloadsize

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func init() {
	caddy.RegisterModule(&PayloadSize{})
	httpcaddyfile.RegisterHandlerDirective("payloadsize", parseCaddyfile)
}

type PayloadSize struct {
	MaxPayloadSize int `json:"max_payload_size,omitempty"` // Optional: Maximum allowed payload size in bytes
	logger         *zap.Logger
}

// UnmarshalCaddyfile Caddyfile syntax:
//
//	payloadsize {
//		max_payload_size <size>
//	}
func (j *PayloadSize) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "max_payload_size":
				if !d.NextArg() {
					return d.ArgErr()
				}
				j.MaxPayloadSize, _ = strconv.Atoi(d.Val())
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	return nil
}

func (j *PayloadSize) Validate() error {
	return nil
}

func (j *PayloadSize) Provision(context caddy.Context) error {
	j.logger = context.Logger()
	return nil
}

// CaddyModule returns the Caddy module information
func (j *PayloadSize) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.payloadsize",
		New: func() caddy.Module { return new(PayloadSize) },
	}
}

func (j *PayloadSize) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	// Read the body (plaintext payload)
	body, err := io.ReadAll(r.Body)

	if err != nil {
		return caddyhttp.Error(http.StatusBadRequest, fmt.Errorf("failed to read request body: %v", err))
	}

	payloadSize := len(body)
	j.logger.Info("recorded payload size", zap.Int("size", payloadSize), zap.String("tenant", "unknown"))

	// Enforce max payload size if specified
	if j.MaxPayloadSize > 0 && payloadSize > j.MaxPayloadSize {
		j.logger.Error("payload size exceeds maximum allowed size", zap.Int("size", payloadSize), zap.Int("max", j.MaxPayloadSize))
		return caddyhttp.Error(http.StatusRequestEntityTooLarge, fmt.Errorf("payload size exceeds maximum allowed size of %d bytes", j.MaxPayloadSize))
	}

	// TODO
	// 1. Determine tenant via JWT claim
	// 2. Store payload size in a database

	// Proceed to the next handler
	return next.ServeHTTP(w, r)
}

// parseCaddyfile will unmarshal tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := &PayloadSize{}
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

var (
	_ caddy.Module                = (*PayloadSize)(nil)
	_ caddyhttp.MiddlewareHandler = (*PayloadSize)(nil)
	_ caddy.Provisioner           = (*PayloadSize)(nil)
	_ caddyfile.Unmarshaler       = (*PayloadSize)(nil)
	_ caddy.Validator             = (*PayloadSize)(nil)
)
