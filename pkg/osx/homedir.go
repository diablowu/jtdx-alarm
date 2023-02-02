package osx

import "os"

func MustUserHomeDir() string {
	if dir, err := os.UserHomeDir(); err != nil {
		panic(err)
	} else {
		return dir
	}
}
