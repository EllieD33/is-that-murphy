package utils

import "strings"

func Canonical(value string) string {
    return strings.ToLower(strings.TrimSpace(value))
}
