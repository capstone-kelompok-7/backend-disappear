package utils

import (
	"github.com/google/uuid"
)

type GeneratorInterface interface {
	GenerateUUID() (string, error)
}

type GeneratorUUID struct{}

func NewGeneratorUUID() GeneratorInterface {
	return &GeneratorUUID{}
}

func (g *GeneratorUUID) GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
