package util

import (
	"os"
	"strings"
)

func ParseEnvFile(fileName string) (result map[string]string) {
	envBuf, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	return ParseEnv(string(envBuf))
}

func ParseEnv(contents string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(contents, "\n") {
		before, after, found := strings.Cut(line, "=")
		if found {
			key := strings.TrimSpace(before)
			val := strings.TrimSpace(after)
			if len(key) > 0 && len(val) > 0 {
				result[key] = val
			}
		}
	}

	return result
}
