package dpdkinit

import (
	"fmt"
	"os"
	"strings"
)

func Parse(key string) (string, error) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			name := strings.TrimPrefix(arg, "--")
			if i+1 >= len(args) {
				return "", fmt.Errorf("no value provided for flag: %s", name)
			}
			return args[i+1], nil
		}

		if strings.HasPrefix(arg, "-") {
			shorthand := strings.TrimPrefix(arg, "-")
			if i+1 >= len(args) {
				return "", fmt.Errorf("no value provided for flag: %s", shorthand)
			}
			return args[i+1], nil
		}
	}

	return "", fmt.Errorf("cant find this key: %s", key)
}
