package config

import "os"

func GetIsolate() string {
	if _, exist := os.LookupEnv("ISOLATE"); !exist {
		return "/usr/local/bin/isolate"
	}

	return os.Getenv("ISOLATE")
}
