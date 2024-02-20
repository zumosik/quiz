package utils

import "fmt"

func WrapErr(msg string, err error) string {
	return fmt.Sprintf("%s: %v", msg, err)
}
