package service

import (
	"fmt"
	"log"
	"time"

	"github.com/Fedasov/Effective-Mobile/internal/model"
	"github.com/Fedasov/Effective-Mobile/internal/repository"
)

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *subscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(req model.SubscriptionCreateRequest) (*model.Subscription, error) {
	log.Printf("Creating subscription for user %s to service %s", req.UserID, req.ServiceName)

	// Преобразование дат из строкового формата
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %v", err)
		}
		endDate = &parsedEndDate
	}

	// Создание модели подписки
	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Сохранение в репозитории
	if err := s.repo.Create(subscription); err != nil {
		log.Printf("Error creating subscription: %v", err)
		return nil, fmt.Errorf("failed to create subscription: %v", err)
	}

	log.Printf("Subscription created successfully with ID: %d", subscription.ID)
	return subscription, nil
}

func (s *subscriptionService) GetByID(id uint32) (*model.Subscription, error) {
	log.Printf("Getting subscription with ID: %d", id)

	subscription, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting subscription %d: %v", id, err)
		return nil, fmt.Errorf("failed to get subscription: %v", err)
	}

	return subscription, nil
}

func (s *subscriptionService) Update(id uint32, req model.SubscriptionCreateRequest) (*model.Subscription, error) {
	log.Printf("Updating subscription with ID: %d", id)

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %v", err)
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %v", err)
		}
		endDate = &parsedEndDate
	}

	existing.ServiceName = req.ServiceName
	existing.Price = req.Price
	existing.UserID = req.UserID
	existing.StartDate = startDate
	existing.EndDate = endDate

	if err := s.repo.Update(existing); err != nil {
		log.Printf("Error updating subscription %d: %v", id, err)
		return nil, fmt.Errorf("failed to update subscription: %v", err)
	}

	log.Printf("Subscription %d updated successfully", id)
	return existing, nil
}

func (s *subscriptionService) Delete(id uint32) error {
	log.Printf("Deleting subscription with ID: %d", id)

	if err := s.repo.Delete(id); err != nil {
		log.Printf("Error deleting subscription %d: %v", id, err)
		return fmt.Errorf("failed to delete subscription: %v", err)
	}

	log.Printf("Subscription %d deleted successfully", id)
	return nil
}

func (s *subscriptionService) List(limit, offset int32) ([]model.Subscription, error) {
	log.Printf("Getting subscriptions list with limit: %d, offset: %d", limit, offset)

	subscriptions, err := s.repo.List(limit, offset)
	if err != nil {
		log.Printf("Error getting subscriptions list: %v", err)
		return nil, fmt.Errorf("failed to get subscriptions list: %v", err)
	}

	log.Printf("Retrieved %d subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (s *subscriptionService) CalculateTotalCost(req model.TotalCostRequest) (int32, error) {
	log.Printf("Calculating total cost for period %s to %s", req.StartDate, req.EndDate)

	startPeriod, err := parseMonthYear(req.StartDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start period: %v", err)
	}

	endPeriod, err := parseMonthYear(req.EndDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end period: %v", err)
	}

	endPeriod = time.Date(endPeriod.Year(), endPeriod.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	total, err := s.repo.CalculateTotalCost(startPeriod, endPeriod, req.UserID, req.ServiceName)
	if err != nil {
		log.Printf("Error calculating total cost: %v", err)
		return 0, fmt.Errorf("failed to calculate total cost: %v", err)
	}

	log.Printf("Total cost calculated: %d", total)
	return total, nil
}

// parseMonthYear преобразует строку формата "MM-YYYY" в time.Time
func parseMonthYear(monthYear string) (time.Time, error) {
	layout := "01-2006"
	date, err := time.Parse(layout, monthYear)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
