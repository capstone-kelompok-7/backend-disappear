package utils

import (
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GeneratorInterface interface {
	GenerateUUID() (string, error)
	GenerateOrderID() (string, error)
}

type Generator struct {
	currentOrderNumber int
	db                 *gorm.DB
}

func NewGeneratorUUID(db *gorm.DB) *Generator {
	return &Generator{
		currentOrderNumber: 0,
		db:                 db,
	}
}

func (g *Generator) GenerateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (g *Generator) GenerateOrderID() (string, error) {
	prefix := "P00PRL"
	g.currentOrderNumber++
	orderNumber := fmt.Sprintf("%03d", g.currentOrderNumber)
	orderID := prefix + orderNumber

	exists, err := g.checkIDExists(orderID)
	if err != nil {
		return "", err
	}
	if exists {
		return g.GenerateOrderID()
	}

	return orderID, nil
}

func (g *Generator) checkIDExists(orderID string) (bool, error) {
	var count int64
	if err := g.db.Model(&entities.OrderModels{}).Where("id_order = ?", orderID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
