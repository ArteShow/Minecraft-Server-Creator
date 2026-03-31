package core

import "strings"

func normalizeHostIP(ip string) string {
	trimmed := strings.TrimSpace(ip)
	lower := strings.ToLower(trimmed)
	if lower == "" {
		return ""
	}

	if lower == "localhost" || lower == "127.0.0.1" || lower == "::1" {
		return "host.docker.internal"
	}

	return trimmed
}
