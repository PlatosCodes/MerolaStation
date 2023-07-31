// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type CollectionTrain struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	TrainID     int64     `json:"train_id"`
	CreatedAt   time.Time `json:"created_at"`
	TimesTraded int64     `json:"times_traded"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type TradeOffer struct {
	ID                  int64     `json:"id"`
	OfferedTrain        int64     `json:"offered_train"`
	OfferedTrainOwner   int64     `json:"offered_train_owner"`
	RequestedTrain      int64     `json:"requested_train"`
	RequestedTrainOwner int64     `json:"requested_train_owner"`
	CreatedAt           time.Time `json:"created_at"`
}

type TradeTransaction struct {
	ID                  int64     `json:"id"`
	OfferedTrain        int64     `json:"offered_train"`
	OfferedTrainOwner   int64     `json:"offered_train_owner"`
	RequestedTrain      int64     `json:"requested_train"`
	RequestedTrainOwner int64     `json:"requested_train_owner"`
	CreatedAt           time.Time `json:"created_at"`
}

type Train struct {
	ID           int64     `json:"id"`
	ModelNumber  string    `json:"model_number"`
	Name         string    `json:"name"`
	Value        int64     `json:"value"`
	CreatedAt    time.Time `json:"created_at"`
	Version      int64     `json:"version"`
	LastEditedAt time.Time `json:"last_edited_at"`
}

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	HashedPassword    []byte    `json:"hashed_password"`
	Activated         bool      `json:"activated"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	Version           int64     `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
}

type WishlistTrain struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TrainID   int64     `json:"train_id"`
	CreatedAt time.Time `json:"created_at"`
}
