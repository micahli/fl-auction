// internal/model/auction.go
package model

import "time"

// AuctionStatus represents the current state of an auction
type AuctionStatus string

const (
	AuctionStatusActive  AuctionStatus = "ACTIVE"
	AuctionStatusEnded   AuctionStatus = "ENDED"
	AuctionStatusPending AuctionStatus = "PENDING"
)

// Auction represents a live auction with all its properties
type Auction struct {
	ID              string        `json:"id"`
	StartingBid     float64       `json:"startingBid"`
	CurrentBid      float64       `json:"currentBid"`
	CurrentWinner   *string       `json:"currentWinner"`
	Duration        int           `json:"duration"`
	ExtendedBidding bool          `json:"extendedBidding"`
	StartTime       time.Time     `json:"startTime"`
	EndTime         time.Time     `json:"endTime"`
	Status          AuctionStatus `json:"status"`
	Bids            []Bid         `json:"bids"`
}

// NextBid returns the minimum next valid bid amount
func (a *Auction) NextBid() float64 {
	return a.CurrentBid + 1.0
}

// TimeRemaining returns the number of seconds remaining in the auction
func (a *Auction) TimeRemaining() int {
	if a.Status != AuctionStatusActive {
		return 0
	}
	
	remaining := a.EndTime.Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	
	return int(remaining.Seconds())
}

// IsActive checks if the auction is currently active
func (a *Auction) IsActive() bool {
	return a.Status == AuctionStatusActive && time.Now().Before(a.EndTime)
}

// HasBids returns true if at least one bid has been placed
func (a *Auction) HasBids() bool {
	return len(a.Bids) > 0
}

// HighestBid returns the highest bid placed, or nil if no bids exist
func (a *Auction) HighestBid() *Bid {
	if !a.HasBids() {
		return nil
	}
	return &a.Bids[len(a.Bids)-1]
}

// ShouldExtend determines if the auction should be extended based on
// the extended bidding rules and current time
func (a *Auction) ShouldExtend() bool {
	if !a.ExtendedBidding {
		return false
	}
	
	timeRemaining := a.EndTime.Sub(time.Now())
	return timeRemaining < 10*time.Second
}