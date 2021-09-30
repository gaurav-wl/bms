package keyProvider

import (
	"github.com/gauravcoco/bms/providers"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

const defaultABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&$"

type KeyProvider struct {
	ShortUID *shortid.Shortid
}


func NewKeyProvider() providers.KeyProvider {
	sid, err := shortid.New(1, defaultABC, 2342)
	if err != nil {
		logrus.Fatal("Error creating new teris key generator err: ", err)
	}

	return &KeyProvider{
		ShortUID: sid,
	}
}

func (k *KeyProvider) GenerateUniqueKey() string {
	return k.ShortUID.MustGenerate()
}
