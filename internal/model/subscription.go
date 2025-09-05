package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint32     `json:"id" example:"1"`
	ServiceName string     `json:"service_name" example:"Yandex Plus"`
	Price       int32      `json:"price" example:"400"`
	UserID      uuid.UUID  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionCreateRequest struct {
	ServiceName string    `json:"service_name" example:"Yandex Plus" validate:"required"`
	Price       int32     `json:"price" example:"400" validate:"required,gt=0"`
	UserID      uuid.UUID `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" validate:"required"`
	StartDate   string    `json:"start_date" example:"07-2025" validate:"required"`
	EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
}

type TotalCostRequest struct {
	StartDate   string     `json:"start_date" example:"01-2025" validate:"required"`
	EndDate     string     `json:"end_date" example:"12-2025" validate:"required"`
	UserID      *uuid.UUID `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName *string    `json:"service_name,omitempty" example:"Yandex Plus"`
}
