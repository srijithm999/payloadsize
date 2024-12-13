package payloadsize_test

import (
	"bytes"
	"testing"

	"github.com/caddyserver/caddy/v2/caddytest"
)

func TestCaddyfilePayloadsize(t *testing.T) {
	// Admin API must be exposed on port 2999 to match what caddytest.Tester does
	config := `
	{
		skip_install_trust
		admin 127.0.0.1:2999
        order payloadsize first
	}

	http://127.0.0.1:12344 {
		bind 127.0.0.1

		payloadsize { 
         max_payload_size 100000	
        }
	    respond 200
	}
	`

	tester := caddytest.NewTester(t)
	tester.InitServer(config, "caddyfile")

	body := bytes.NewBufferString("foo")

	tester.AssertPostResponseBody("http://127.0.0.1:12344", []string{"Content-Type: application/json"}, body, 200, "")
}
