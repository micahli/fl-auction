package model

// AuctionEventType represents the type of auction event
type AuctionEventType string

const (
	EventAuctionStarted AuctionEventType = "AUCTION_STARTED"
	EventBidPlaced      AuctionEventType = "BID_PLACED"
	EventAuctionEnded   AuctionEventType = "AUCTION_ENDED"
)

// AuctionEvent represents an event that occurred in the auction system
type AuctionEvent struct {
	Type    AuctionEventType `json:"type"`
	Auction *Auction         `json:"auction,omitempty"`
	Bid     *Bid             `json:"bid,omitempty"`
	Error   *string          `json:"error,omitempty"`
}

// NewAuctionStartedEvent creates an event for when an auction starts
func NewAuctionStartedEvent(auction *Auction) *AuctionEvent {
	return &AuctionEvent{
		Type:    EventAuctionStarted,
		Auction: auction,
	}
}

// NewBidPlacedEvent creates an event for when a bid is placed
func NewBidPlacedEvent(auction *Auction, bid *Bid) *AuctionEvent {
	return &AuctionEvent{
		Type:    EventBidPlaced,
		Auction: auction,
		Bid:     bid,
	}
}

// NewAuctionEndedEvent creates an event for when an auction ends
func NewAuctionEndedEvent(auction *Auction) *AuctionEvent {
	return &AuctionEvent{
		Type:    EventAuctionEnded,
		Auction: auction,
	}
}

// NewErrorEvent creates an event for when an error occurs
func NewErrorEvent(errMsg string) *AuctionEvent {
	return &AuctionEvent{
		Error: &errMsg,
	}
}

// IsError checks if this event represents an error
func (e *AuctionEvent) IsError() bool {
	return e.Error != nil
}
