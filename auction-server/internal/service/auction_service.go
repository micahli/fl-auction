package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/micahli/fl-auction/auction-server/internal/model"
	"github.com/micahli/fl-auction/auction-server/internal/store"
)

// AuctionService handles auction business logic
type AuctionService struct {
	store          *store.AuctionStore
	validationRule *model.ValidationRules
	timerMutex     sync.Mutex
}

// NewAuctionService creates a new auction service
func NewAuctionService(store *store.AuctionStore) *AuctionService {
	return &AuctionService{
		store:          store,
		validationRule: model.DefaultValidationRules(),
	}
}

// CreateAuction creates and starts a new auction
func (s *AuctionService) CreateAuction(ctx context.Context, startingBid float64, duration int, extendedBidding bool) (*model.Auction, error) {
	s.timerMutex.Lock()
	defer s.timerMutex.Unlock()

	// Validate there's no active auction
	current := s.store.GetCurrentAuction()
	if current != nil && current.Status == model.AuctionStatusActive {
		return nil, model.ErrAuctionAlreadyActive
	}

	// Validate starting bid
	if err := s.validationRule.ValidateStartingBid(startingBid); err != nil {
		return nil, err
	}

	// Validate duration
	if duration <= 0 {
		duration = 30 // Default duration
	}
	if err := s.validationRule.ValidateDuration(duration); err != nil {
		return nil, err
	}

	// Create auction
	now := time.Now()
	auction := &model.Auction{
		ID:              fmt.Sprintf("auction-%d", now.UnixNano()),
		StartingBid:     startingBid,
		CurrentBid:      startingBid,
		CurrentWinner:   nil,
		Duration:        duration,
		ExtendedBidding: extendedBidding,
		StartTime:       now,
		EndTime:         now.Add(time.Duration(duration) * time.Second),
		Status:          model.AuctionStatusActive,
		Bids:            []model.Bid{},
	}

	s.store.SetCurrentAuction(auction)

	// Broadcast auction started event
	s.store.Broadcast(model.NewAuctionStartedEvent(auction))

	// Start countdown timer
	go s.startCountdown(auction)

	return auction, nil
}

// PlaceBid attempts to place a bid on the current auction
func (s *AuctionService) PlaceBid(ctx context.Context, userID string, amount float64) (*model.Bid, error) {
	s.timerMutex.Lock()
	defer s.timerMutex.Unlock()

	auction := s.store.GetCurrentAuction()
	if auction == nil || auction.Status != model.AuctionStatusActive {
		return nil, model.ErrNoActiveAuction
	}

	now := time.Now()

	// Check if bid is too late
	if now.After(auction.EndTime) {
		return nil, model.NewBidTooLateError()
	}

	// Validate bid amount
	if err := s.validationRule.ValidateBidAmount(amount, auction.CurrentBid); err != nil {
		if err == model.ErrBidTooLow {
			return nil, model.NewBidTooLowError(auction.CurrentBid, amount)
		}
		return nil, err
	}

	// Create bid
	bid := &model.Bid{
		ID:        fmt.Sprintf("bid-%d", s.store.GetNextBidID()),
		AuctionID: auction.ID,
		UserID:    userID,
		Amount:    amount,
		Timestamp: now,
	}

	// Add bid to auction
	if err := s.store.AddBid(bid); err != nil {
		return nil, err
	}

	// Handle extended bidding
	if s.validationRule.ShouldExtendAuction(auction.EndTime, auction.ExtendedBidding) {
		auction.EndTime = s.validationRule.CalculateExtendedEndTime(now)
	}

	// Broadcast bid placed event
	s.store.Broadcast(model.NewBidPlacedEvent(auction, bid))

	return bid, nil
}

// GetCurrentAuction returns the current auction
func (s *AuctionService) GetCurrentAuction() *model.Auction {
	return s.store.GetCurrentAuction()
}

// GetNextBid returns the minimum next valid bid
func (s *AuctionService) GetNextBid() float64 {
	auction := s.store.GetCurrentAuction()
	if auction == nil {
		return 0
	}
	return s.validationRule.CalculateNextMinimumBid(auction.CurrentBid)
}

// GetTimeRemaining returns seconds remaining in the auction
func (s *AuctionService) GetTimeRemaining() int {
	auction := s.store.GetCurrentAuction()
	if auction == nil || auction.Status != model.AuctionStatusActive {
		return 0
	}

	return auction.TimeRemaining()
}

// startCountdown runs a countdown timer for the auction
func (s *AuctionService) startCountdown(auction *model.Auction) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		current := s.store.GetCurrentAuction()
		if current == nil || current.ID != auction.ID {
			return
		}

		if time.Now().After(current.EndTime) {
			s.endAuction(current)
			return
		}
	}
}

// endAuction marks an auction as ended and broadcasts the event
func (s *AuctionService) endAuction(auction *model.Auction) {
	s.timerMutex.Lock()
	defer s.timerMutex.Unlock()

	auction.Status = model.AuctionStatusEnded
	s.store.SetCurrentAuction(auction)

	// Broadcast auction ended event
	s.store.Broadcast(model.NewAuctionEndedEvent(auction))
}

// Subscribe creates a new event subscription
func (s *AuctionService) Subscribe(id string) chan *model.AuctionEvent {
	return s.store.Subscribe(id)
}

// Unsubscribe removes an event subscription
func (s *AuctionService) Unsubscribe(id string) {
	s.store.Unsubscribe(id)
}

// GetSubscriberCount returns the number of active subscribers
func (s *AuctionService) GetSubscriberCount() int {
	return s.store.GetSubscriberCount()
}
