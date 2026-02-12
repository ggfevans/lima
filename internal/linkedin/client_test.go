package linkedin

import (
	"strings"
	"testing"
)

func TestNormalizeCookieInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantLiAt bool
	}{
		{
			name:     "full cookie header",
			input:    `bcookie="v=2&abc123"; JSESSIONID="ajax:123456"; li_at=AQEFtoken123; liap=true`,
			wantErr:  false,
			wantLiAt: true,
		},
		{
			name:     "just li_at assignment",
			input:    "li_at=AQEFtoken123value",
			wantErr:  false,
			wantLiAt: true,
		},
		{
			name:     "bare token value",
			input:    "AQEFtoken123value",
			wantErr:  false,
			wantLiAt: true,
		},
		{
			name:     "empty",
			input:    "",
			wantErr:  true,
			wantLiAt: false,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			wantErr:  true,
			wantLiAt: false,
		},
		{
			name:     "full header with quoted values and special chars",
			input:    `bcookie="v=2&4c05e43e-71e7-4e43-8dbf-803ea2503aa5"; JSESSIONID="ajax:1317978622907181656"; li_at=AQEFAHMBAAAAABvundIAAAGcSaDdcAAAAZxtrWFw+test/value=; liap=true`,
			wantErr:  false,
			wantLiAt: true,
		},
		{
			name:     "li_at with slashes and plus signs",
			input:    "li_at=AQEFtoken/with+special/chars=ending==",
			wantErr:  false,
			wantLiAt: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeCookieInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeCookieInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantLiAt && !strings.Contains(result, "li_at=") {
				t.Errorf("normalizeCookieInput() result %q does not contain li_at=", result)
			}
		})
	}
}

func TestExtractCookieValue(t *testing.T) {
	header := `bcookie="v=2&abc"; JSESSIONID="ajax:12345"; li_at=AQEFtoken; liap=true`

	liAt := extractCookieValue(header, "li_at")
	if liAt != "AQEFtoken" {
		t.Errorf("expected AQEFtoken, got %q", liAt)
	}

	jsess := extractCookieValue(header, "JSESSIONID")
	if jsess != `"ajax:12345"` {
		t.Errorf("expected \"ajax:12345\", got %q", jsess)
	}

	missing := extractCookieValue(header, "nonexistent")
	if missing != "" {
		t.Errorf("expected empty, got %q", missing)
	}
}
