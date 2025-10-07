package model

import "time"

// Bid represents a single bid placed by a user
type Bid struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auctionId"`
	UserID    string    `json:"userId"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// IsHigherThan checks if this bid amount is higher than the given amount
func (b *Bid) IsHigherThan(amount float64) bool {
	return b.Amount > amount
}

// IsPlacedBy checks if this bid was placed by the given user
func (b *Bid) IsPlacedBy(userID string) bool {
	return b.UserID == userID
}

// Age returns the duration since this bid was placed
func (b *Bid) Age() time.Duration {
	return time.Since(b.Timestamp)
}
