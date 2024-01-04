package util

import "os"

func CreateTmpFile() (string, error) {
	f, err := os.CreateTemp("", "nginx.conf")
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Name(), nil
}
