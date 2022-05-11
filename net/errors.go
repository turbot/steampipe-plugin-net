package net

func shouldRetryError(err error) bool {
	return err != nil
}
