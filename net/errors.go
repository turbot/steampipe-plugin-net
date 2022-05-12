package net

import "strings"

func shouldRetryError(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			return true
		}
	}
	return false
}
