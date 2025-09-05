package service

import (
	"github.com/Fedasov/Effective-Mobile/internal/model"
)

type SubscriptionService interface {
	Create(req model.SubscriptionCreateRequest) (*model.Subscription, error)
	GetByID(id uint32) (*model.Subscription, error)
	Update(id uint32, req model.SubscriptionCreateRequest) (*model.Subscription, error)
	Delete(id uint32) error
	List(limit, offset int32) ([]model.Subscription, error)
	CalculateTotalCost(req model.TotalCostRequest) (int32, error)
}
