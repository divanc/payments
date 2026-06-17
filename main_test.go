package main

import (
	"testing"

	"github.com/divanc/payments/gateway"
)

func TestNewGateway(t *testing.T) {
	cases := []struct {
		name      string
		gw        string
		key       string
		wantType  any
		wantError bool
	}{
		{"default empty", "", "", (*gateway.Fake)(nil), false},
		{"explicit fake", "fake", "", (*gateway.Fake)(nil), false},
		{"stripe test key", "stripe", "sk_test_abc123", (*gateway.Stripe)(nil), false},
		{"stripe missing key", "stripe", "", nil, true},
		{"stripe live key", "stripe", "sk_live_abc123", nil, true},
		{"unknown", "bogus", "", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gw, err := newGateway(tc.gw, tc.key)
			if tc.wantError {
				if err == nil {
					t.Fatalf("newGateway(%q, …) = ok, want error", tc.gw)
				}
				return
			}
			if err != nil {
				t.Fatalf("newGateway(%q, …) = err %v, want ok", tc.gw, err)
			}
			switch tc.wantType.(type) {
			case *gateway.Fake:
				if _, ok := gw.(*gateway.Fake); !ok {
					t.Fatalf("newGateway(%q, …) = %T, want *gateway.Fake", tc.gw, gw)
				}
			case *gateway.Stripe:
				if _, ok := gw.(*gateway.Stripe); !ok {
					t.Fatalf("newGateway(%q, …) = %T, want *gateway.Stripe", tc.gw, gw)
				}
			}
		})
	}
}
