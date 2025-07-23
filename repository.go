package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *Subscription) error
	Read(id uint) (*Subscription, error)
	Update(subscription *Subscription) error
	Delete(id uint) error
	List() ([]Subscription, error)
	Filter(startDate time.Time, endDate time.Time, userID uuid.UUID, serviceName string) ([]Subscription, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) Read(id uint) (*Subscription, error) {
	var subscription Subscription
	if err := r.db.First(&subscription, id).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) Update(subscription *Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *subscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&Subscription{}, id).Error
}

func (r *subscriptionRepository) List() ([]Subscription, error) {
	var subs []Subscription
	if err := r.db.Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *subscriptionRepository) Filter(startDate time.Time, endDate time.Time, userID uuid.UUID, serviceName string) ([]Subscription, error) {
	var subs []Subscription

	query := r.db.Model(&Subscription{})

	// We find subscriptions even if they fit in the time range partially
	// If end date is not set for subscription, it is considered endless and therefore not taken into account
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("start_date <= ? AND end_date >= ?", endDate, startDate)
	} else if !startDate.IsZero() {
		query = query.Where("start_date >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("end_date <= ?", endDate)
	}

	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	}

	return subs, nil
}
