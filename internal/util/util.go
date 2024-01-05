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

var dotlocalPath *string

func GetDotlocalPath() string {
	if dotlocalPath == nil {
		home := os.Getenv("HOME")
		dir := path.Join(home, ".dotlocal")
		dotlocalPath = &dir

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	return *dotlocalPath
}

func GetApiSocketPath() string {
	return path.Join(GetDotlocalPath(), "api.sock")
}

func GetPidPath() string {
	return path.Join(GetDotlocalPath(), "pid")
}

func GetPreferencesPath() string {
	return path.Join(GetDotlocalPath(), "preferences.json")
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
