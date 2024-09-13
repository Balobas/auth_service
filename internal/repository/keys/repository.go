package repositoryKeys

import (
	"errors"
	"os"
)

type KeysRepository struct{}

func New() *KeysRepository {
	return &KeysRepository{}
}

func (r *KeysRepository) GetPrivateKey() ([]byte, error) {
	pk := os.Getenv("PK")
	if len(pk) == 0 {
		return nil, errors.New("failed to get private key")
	}
	return []byte(pk), nil
}
