package utils

import "fmt"

func ToString(err error) string {
	return fmt.Sprintf("%+v", err)
}
