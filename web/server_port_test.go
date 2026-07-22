package main

import "testing"

func TestResolveListenAddress(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		set     bool
		want    string
		wantErr bool
	}{
		{name: "unset", want: ":8080"},
		{name: "empty", value: "", set: true, want: ":8080"},
		{name: "custom", value: "18080", set: true, want: ":18080"},
		{name: "zero", value: "0", set: true, wantErr: true},
		{name: "too large", value: "65536", set: true, wantErr: true},
		{name: "text", value: "http", set: true, wantErr: true},
		{name: "whitespace", value: " 18080", set: true, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.set {
				t.Setenv("TDX_API_PORT", test.value)
			}
			got, err := resolveListenAddress()
			if test.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
