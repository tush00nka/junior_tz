package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint       `json:"id" gorm:"type:uuid"`
	ServiceName string     `json:"service_name"`
	Price       uint       `json:"price"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func (s *Subscription) UnmarshalJSON(data []byte) error {
	type Alias Subscription
	aux := &struct {
		StartDate string  `json:"start_date"`
		EndDate   *string `json:"end_date,omitempty"` // Теперь *string
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Парсинг StartDate (обязательное поле)
	startDate, err := time.Parse("01-2006", aux.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format: expected MM-YYYY")
	}
	s.StartDate = startDate

	// Парсинг EndDate (необязательное поле)
	if aux.EndDate != nil {
		endDate, err := time.Parse("01-2006", *aux.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date format: expected MM-YYYY")
		}
		s.EndDate = &endDate
	} else {
		s.EndDate = nil
	}

	return nil
}
