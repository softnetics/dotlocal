package util

import (
	"net"
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

func FindAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	if err != nil {
		return 0, err
	}
	return port, nil
}
