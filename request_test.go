package server

import "testing"

func TestParseRequestLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		line       string
		wantMethod string
		wantURI    string
		wantProto  string
		wantOK     bool
	}{
		{
			name:       "валидный GET запрос",
			line:       "GET /hello HTTP/1.1",
			wantMethod: "GET",
			wantURI:    "/hello",
			wantProto:  "HTTP/1.1",
			wantOK:     true,
		},
		{
			name:       "валидный POST запрос",
			line:       "POST /api/user HTTP/1.1",
			wantMethod: "POST",
			wantURI:    "/api/user",
			wantProto:  "HTTP/1.1",
			wantOK:     true,
		},
		{
			name: "пустая строка",
			line: "",
		},
		{
			name:       "один элемент (нет версии)",
			line:       "GET /hello",
			wantMethod: "",
			wantURI:    "",
			wantProto:  "",
			wantOK:     false,
		},
		{
			name:       "один элемент вместо трёх",
			line:       "GET",
			wantMethod: "",
			wantURI:    "",
			wantProto:  "",
			wantOK:     false,
		},
		{
			name:       "четыре элемента",
			line:       "GET / HTTP/1.1 EXTRA",
			wantMethod: "GET",
			wantURI:    "/",
			wantProto:  "HTTP/1.1 EXTRA",
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotURI, gotProto, gotOK := parseRequestLine(tt.line)

			if gotMethod != tt.wantMethod {
				t.Errorf("Method: got %q, want %q", gotMethod, tt.wantMethod)
			}
			if gotURI != tt.wantURI {
				t.Errorf("URI: got %q, want %q", gotURI, tt.wantURI)
			}
			if gotProto != tt.wantProto {
				t.Errorf("Proto: got %q, want %q", gotProto, tt.wantProto)
			}
			if gotOK != tt.wantOK {
				t.Errorf("OK: got %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}
