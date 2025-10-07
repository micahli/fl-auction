package model

import "errors"

// Common errors used throughout the auction system
var (
	ErrAuctionAlreadyActive = errors.New("an auction is already active")
	ErrNoActiveAuction      = errors.New("no active auction")
	ErrBidTooLow            = errors.New("bid too low")
	ErrBidTooLate           = errors.New("bid too late")
	ErrInvalidBidAmount     = errors.New("invalid bid amount")
	ErrInvalidDuration      = errors.New("invalid auction duration")
	ErrInvalidStartingBid   = errors.New("invalid starting bid")
	ErrAuctionNotFound      = errors.New("auction not found")
)

// BidError represents a bid-specific error with context
type BidError struct {
	Err           error
	CurrentBid    float64
	AttemptedBid  float64
	TimeRemaining int
}

func (e *BidError) Error() string {
	return e.Err.Error()
}

func (e *BidError) Unwrap() error {
	return e.Err
}

// NewBidTooLowError creates a detailed bid too low error
func NewBidTooLowError(currentBid, attemptedBid float64) *BidError {
	return &BidError{
		Err:          ErrBidTooLow,
		CurrentBid:   currentBid,
		AttemptedBid: attemptedBid,
	}
}

// NewBidTooLateError creates a detailed bid too late error
func NewBidTooLateError() *BidError {
	return &BidError{
		Err: ErrBidTooLate,
	}
}
