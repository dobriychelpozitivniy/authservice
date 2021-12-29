package service

import (
	"testing"
)



func TestTokenManagerService_ParseAccessToken(t *testing.T) {
	tests := []struct {
		name       string
		singingKey string
		token      string
		wantErr    bool
	}{
		{
			name:       "vaild",
			singingKey: "asdf",
			token:      "",
			wantErr:    false,
		},
		{
			name:       "invalid",
			singingKey: "asdf",
			token:      "invalid",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TokenManagerService{
				signingKey: tt.singingKey,
			}
			if tt.token == "" {
				tt.token, _ = s.GenerateAccessToken("testuser")
			}
			_, err := s.ParseAccessToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTokenManagerService_ParseRefreshToken(t *testing.T) {
	tests := []struct {
		name       string
		singingKey string
		token      string
		wantErr    bool
	}{
		{
			name:       "vaild",
			singingKey: "asdf",
			token:      "",
			wantErr:    false,
		},
		{
			name:       "invalid",
			singingKey: "asdf",
			token:      "invalid",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TokenManagerService{
				signingKey: tt.singingKey,
			}
			if tt.token == "" {
				tt.token, _ = s.GenerateRefreshToken("testuser")
			}
			_, err := s.ParseRefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
