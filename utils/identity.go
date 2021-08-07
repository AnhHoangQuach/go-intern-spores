package utils

import (
	"math/big"
	"strings"

	"github.com/google/uuid"
)

type IIdentity interface {
	NewString() string
	NewBigIntString() string
}

type Guuid struct {
}

func NewGuuid() Guuid {
	return Guuid{}
}
func (g Guuid) NewString() string {
	return uuid.NewString()
}
func (g Guuid) NewBigIntString() string {
	id := uuid.NewString()
	var i big.Int
	i.SetString(strings.Replace(id, "-", "", 4), 16)
	return i.String()
}
