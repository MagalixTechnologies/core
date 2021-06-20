package cache

import (
	"net/http"
	"strings"
)

type URLMatch struct {
	Path   string
	Method string
}

func (m *URLMatch) MatchRequest(req *http.Request) bool {
	if m.Method != req.Method {
		return false
	}
	if strings.Contains(m.Path, ":arg") {
		parts := strings.Split(m.Path, "/")
		requestParts := strings.Split(req.URL.Path, "/")
		if len(parts) != len(requestParts) {
			return false
		}
		for i := range parts {
			if parts[i] == ":arg" {
				continue
			} else {
				if parts[i] != requestParts[i] {
					return false
				}
			}

		}
		return true
	} else {
		return req.URL.Path == m.Path
	}
}
