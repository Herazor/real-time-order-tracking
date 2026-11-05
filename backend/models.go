package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Role      string    `gorm:"not null"`
	CreatedAt time.Time
}

type Order struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ClientID    uuid.UUID `gorm:"type:uuid;not null"`
	OfficerID   *uuid.UUID
	ServiceName string    `gorm:"not null"`
	Description string
	Status      string    `gorm:"not null;default:PENDING"`
	Version     int       `gorm:"default:1"`
	LastSequence int      `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderEvent struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null"`
	Sequence       int       `gorm:"not null"`
	OldStatus      string
	NewStatus      string    `gorm:"not null"`
	ActorID        uuid.UUID `gorm:"type:uuid;not null"`
	IdempotencyKey *string `gorm:"uniqueIndex"`
	Note           string
	CreatedAt      time.Time
}

type CancellationRequest struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    OrderID    uuid.UUID `gorm:"type:uuid;not null"`
    RequestedBy uuid.UUID `gorm:"type:uuid;not null"`
    Reason     string    `gorm:"not null"`
    Status     string    `gorm:"not null;default:PENDING"`
    ReviewedBy *uuid.UUID
    CreatedAt  time.Time
    UpdatedAt  time.Time
}