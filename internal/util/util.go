package util

import (
	"os"
	"path"
)

func CreateTmpFile() (string, error) {
	f, err := os.CreateTemp("", "nginx.conf")
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Name(), nil
}

func GetApiSocketPath() string {
	home := os.Getenv("HOME")

	err := os.MkdirAll(path.Join(home, ".dotlocal"), 0755)
	if err != nil {
		panic(err)
	}

	return path.Join(home, ".dotlocal", "api.sock")
}
