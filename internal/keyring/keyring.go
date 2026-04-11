package keyring

import (
	"github.com/zalando/go-keyring"
)

const (
	serviceName = "sub-muse"
)

func SavePassword(username, password string) error {
	return keyring.Set(serviceName, username, password)
}

func GetPassword(username string) (string, error) {
	return keyring.Get(serviceName, username)
}

func DeletePassword(username string) error {
	return keyring.Delete(serviceName, username)
}

func IsAvailable() bool {
	_ = keyring.Set("sub-muse", "test", "test")
	err := keyring.Delete("sub-muse", "test")
	return err == nil
}
