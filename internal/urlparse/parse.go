// Package urlparse extracts org and project name from git remote URLs.
// It handles three forms:
//
//   - HTTPS:  https://github.com/org/project.git
//   - SSH:    git@github.com:org/project.git
//   - Local:  bare name (no slashes, no URL scheme) → org="local"
package urlparse

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Result holds the parsed components of a git URL.
type Result struct {
	Org      string // e.g. "myorg" or "local"
	Project  string // e.g. "myrepo"
	CloneURL string // original URL, unchanged (empty for local names)
}

var sshPattern = regexp.MustCompile(`^git@[^:]+:([^/]+)/([^/]+?)(?:\.git)?$`)

// Parse inspects input and returns a Result. Returns an error only when the
// input looks like a URL but cannot be parsed.
func Parse(input string) (Result, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{}, fmt.Errorf("input is empty")
	}

	// SSH URL: git@host:org/project[.git]
	if m := sshPattern.FindStringSubmatch(input); m != nil {
		return Result{
			Org:      m[1],
			Project:  m[2],
			CloneURL: input,
		}, nil
	}

	// HTTPS URL
	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://") {
		u, err := url.Parse(input)
		if err != nil {
			return Result{}, fmt.Errorf("invalid URL: %w", err)
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) < 2 {
			return Result{}, fmt.Errorf("URL path must contain org/project: %s", input)
		}
		org := parts[0]
		project := strings.TrimSuffix(parts[1], ".git")
		if org == "" || project == "" {
			return Result{}, fmt.Errorf("could not extract org/project from URL: %s", input)
		}
		return Result{
			Org:      org,
			Project:  project,
			CloneURL: input,
		}, nil
	}

	// Local name (no URL scheme, no colon-path)
	return Result{
		Org:      "local",
		Project:  input,
		CloneURL: "",
	}, nil
}
