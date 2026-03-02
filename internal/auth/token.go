package auth

import (
	"encoding/json"

	"github.com/99designs/keyring"
	"github.com/voska/amadeus-cli/internal/config"
)

const serviceName = "amadeus-cli"

func openKeyring() (keyring.Keyring, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}
	return keyring.Open(keyring.Config{
		ServiceName:      serviceName,
		FileDir:          dir + "/tokens",
		FilePasswordFunc: func(string) (string, error) { return "", nil },
	})
}

func StoreToken(env string, token *Token) error {
	kr, err := openKeyring()
	if err != nil {
		return err
	}
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return kr.Set(keyring.Item{
		Key:  env,
		Data: data,
	})
}

func LoadToken(env string) (*Token, error) {
	kr, err := openKeyring()
	if err != nil {
		return nil, err
	}
	item, err := kr.Get(env)
	if err != nil {
		return nil, err
	}
	var tok Token
	if err := json.Unmarshal(item.Data, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func DeleteToken(env string) error {
	kr, err := openKeyring()
	if err != nil {
		return err
	}
	return kr.Remove(env)
}
