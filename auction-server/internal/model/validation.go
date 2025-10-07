package model

import "time"

// ValidationRules contains configuration for auction validation
type ValidationRules struct {
	MinStartingBid     float64
	MaxStartingBid     float64
	MinDuration        int
	MaxDuration        int
	MinBidIncrement    float64
	ExtensionThreshold time.Duration
	ExtensionDuration  time.Duration
}

// DefaultValidationRules returns the default validation rules
func DefaultValidationRules() *ValidationRules {
	return &ValidationRules{
		MinStartingBid:     1.0,
		MaxStartingBid:     1000000.0,
		MinDuration:        10,
		MaxDuration:        3600,
		MinBidIncrement:    1.0,
		ExtensionThreshold: 10 * time.Second,
		ExtensionDuration:  10 * time.Second,
	}
}

// ValidateStartingBid checks if the starting bid is valid
func (vr *ValidationRules) ValidateStartingBid(amount float64) error {
	if amount < vr.MinStartingBid {
		return ErrInvalidStartingBid
	}
	if amount > vr.MaxStartingBid {
		return ErrInvalidStartingBid
	}
	return nil
}

// ValidateDuration checks if the auction duration is valid
func (vr *ValidationRules) ValidateDuration(duration int) error {
	if duration < vr.MinDuration {
		return ErrInvalidDuration
	}
	if duration > vr.MaxDuration {
		return ErrInvalidDuration
	}
	return nil
}

// ValidateBidAmount checks if a bid amount is valid
func (vr *ValidationRules) ValidateBidAmount(amount, currentBid float64) error {
	if amount <= 0 {
		return ErrInvalidBidAmount
	}
	if amount <= currentBid {
		return ErrBidTooLow
	}
	// Optionally enforce minimum increment
	// if amount < currentBid + vr.MinBidIncrement {
	//     return ErrBidTooLow
	// }
	return nil
}

// CalculateNextMinimumBid calculates the next valid minimum bid
func (vr *ValidationRules) CalculateNextMinimumBid(currentBid float64) float64 {
	return currentBid + vr.MinBidIncrement
}

// ShouldExtendAuction determines if an auction should be extended
func (vr *ValidationRules) ShouldExtendAuction(endTime time.Time, extendedBiddingEnabled bool) bool {
	if !extendedBiddingEnabled {
		return false
	}

	timeRemaining := time.Until(endTime)
	return timeRemaining > 0 && timeRemaining < vr.ExtensionThreshold
}

// CalculateExtendedEndTime calculates the new end time after extension
func (vr *ValidationRules) CalculateExtendedEndTime(currentTime time.Time) time.Time {
	return currentTime.Add(vr.ExtensionDuration)
}
