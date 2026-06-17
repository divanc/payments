package gateway

import "testing"

func TestNewStripeRejectsNonTestKeys(t *testing.T) {
	cases := []struct {
		name string
		key  string
		ok   bool
	}{
		{"empty", "", false},
		{"live", "sk_live_abc123", false},
		{"publishable", "pk_test_abc123", false},
		{"junk", "not-a-key", false},
		{"test", "sk_test_abc123", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gw, err := NewStripe(tc.key)
			if tc.ok {
				if err != nil {
					t.Fatalf("NewStripe(%q) = err %v, want ok", tc.key, err)
				}
				if gw == nil {
					t.Fatalf("NewStripe(%q) = nil gateway, want non-nil", tc.key)
				}
				return
			}
			if err == nil {
				t.Fatalf("NewStripe(%q) = ok, want error", tc.key)
			}
		})
	}
}
