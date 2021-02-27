package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Day struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	MonthID   uuid.UUID `json:"month_id"`
	DayOfWeek string    `json:"day_od_week"`
	Hours     float64   `json:"hours"`
	Go        time.Time `json:"go_time"`
	Come      time.Time `json:"come_time"`
}

type Month struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID   int       `json:"user_id"`
	Days     []Day     `json:"weeks" gorm:"foreignKey:MonthID"`
	SumHours int       `json:"sum_hours"`
	Name     string    `json:"name_of_month"`
}
