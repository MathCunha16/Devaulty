package model

import "time"

type BaseEntity struct {
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}
