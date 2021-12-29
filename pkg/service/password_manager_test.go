package service

import (
	"testing"
)

func TestPasswordManagerService_Check(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				password: "wwww",
				hash:     "$2a$12$tc548.1ls5q7Enkgsj5ivuP0LRU1ATyp0TaAWaSWFveZZd59TmeZm",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				password: "wwww",
				hash:     "invalid",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PasswordManagerService{}
			err := m.Check(tt.args.password, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPasswordManagerService_Hash(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "valid",
			args:    args{password: "wwww"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PasswordManagerService{}
			hash, err := m.Hash(tt.args.password)
			err = m.Check(tt.args.password, hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
