package repository

import (
	"time"

	"github.com/Fedasov/Effective-Mobile/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id uint32) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id uint32) error
	List(limit, offset int32) ([]model.Subscription, error)
	CalculateTotalCost(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int32, error)
}
