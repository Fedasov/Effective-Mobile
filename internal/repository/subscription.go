package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Fedasov/Effective-Mobile/internal/model"

	"github.com/google/uuid"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *subscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(sub *model.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
}

func (r *subscriptionRepository) GetByID(id uint32) (*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date 
	          FROM subscriptions WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var sub model.Subscription
	var endDate sql.NullTime

	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get subscription: %v", err)
	}

	if endDate.Valid {
		sub.EndDate = &endDate.Time
	}

	return &sub, nil
}

func (r *subscriptionRepository) Update(sub *model.Subscription) error {
	query := `UPDATE subscriptions 
	          SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 
	          WHERE id = $6`

	result, err := r.db.Exec(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription with ID %d not found", sub.ID)
	}

	log.Printf("Updated subscription with ID: %d", sub.ID)
	return nil
}

func (r *subscriptionRepository) Delete(id uint32) error {
	query := "DELETE FROM subscriptions WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription with ID %d not found", id)
	}

	log.Printf("Deleted subscription with ID: %d", id)
	return nil
}

func (r *subscriptionRepository) List(limit, offset int32) ([]model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date 
	          FROM subscriptions 
	          ORDER BY id 
	          LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions list: %v", err)
	}
	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {
		var sub model.Subscription
		var endDate sql.NullTime

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %v", err)
		}

		if endDate.Valid {
			sub.EndDate = &endDate.Time
		}

		subscriptions = append(subscriptions, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subscriptions: %v", err)
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) CalculateTotalCost(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int32, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions 
	          WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)`
	args := []interface{}{endDate, startDate}

	paramIndex := 3
	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", paramIndex)
		args = append(args, *userID)
		paramIndex++
	}
	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", paramIndex)
		args = append(args, *serviceName)
	}

	var total int32
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %v", err)
	}

	return total, nil
}
