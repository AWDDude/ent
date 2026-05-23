package urlparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantOrg  string
		wantProj string
		wantURL  string
		wantErr  bool
	}{
		{
			name:     "SSH with .git",
			input:    "git@github.com:myorg/myrepo.git",
			wantOrg:  "myorg",
			wantProj: "myrepo",
			wantURL:  "git@github.com:myorg/myrepo.git",
		},
		{
			name:     "SSH without .git",
			input:    "git@github.com:myorg/myrepo",
			wantOrg:  "myorg",
			wantProj: "myrepo",
			wantURL:  "git@github.com:myorg/myrepo",
		},
		{
			name:     "SSH with subdomain",
			input:    "git@gitlab.company.com:team/tool.git",
			wantOrg:  "team",
			wantProj: "tool",
			wantURL:  "git@gitlab.company.com:team/tool.git",
		},
		{
			name:     "HTTPS with .git",
			input:    "https://github.com/myorg/myrepo.git",
			wantOrg:  "myorg",
			wantProj: "myrepo",
			wantURL:  "https://github.com/myorg/myrepo.git",
		},
		{
			name:     "HTTPS without .git",
			input:    "https://github.com/myorg/myrepo",
			wantOrg:  "myorg",
			wantProj: "myrepo",
			wantURL:  "https://github.com/myorg/myrepo",
		},
		{
			name:     "HTTP (non-TLS)",
			input:    "http://internal.example.com/team/service.git",
			wantOrg:  "team",
			wantProj: "service",
			wantURL:  "http://internal.example.com/team/service.git",
		},
		{
			name:     "local name only",
			input:    "myapp",
			wantOrg:  "local",
			wantProj: "myapp",
			wantURL:  "",
		},
		{
			name:     "local name with hyphens",
			input:    "my-cool-app",
			wantOrg:  "local",
			wantProj: "my-cool-app",
			wantURL:  "",
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "HTTPS missing project",
			input:   "https://github.com/onlyorg",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantOrg, got.Org)
			assert.Equal(t, tt.wantProj, got.Project)
			assert.Equal(t, tt.wantURL, got.CloneURL)
		})
	}
}
