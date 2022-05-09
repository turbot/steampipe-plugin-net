package net

func shouldRetryError(err error) bool {
	if err != nil {
		return true
	}
	return false
}
