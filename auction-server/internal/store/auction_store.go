package store

import (
	"sync"

	"github.com/micahli/fl-auction/auction-server/internal/model"
)

// AuctionStore manages auction state and subscriptions
type AuctionStore struct {
	mu             sync.RWMutex
	currentAuction *model.Auction
	subscribers    map[string]chan *model.AuctionEvent
	nextBidID      int
}

// NewAuctionStore creates a new auction store
func NewAuctionStore() *AuctionStore {
	return &AuctionStore{
		subscribers: make(map[string]chan *model.AuctionEvent),
		nextBidID:   1,
	}
}

// GetCurrentAuction returns the current active auction
func (s *AuctionStore) GetCurrentAuction() *model.Auction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentAuction
}

// SetCurrentAuction sets the current auction
func (s *AuctionStore) SetCurrentAuction(auction *model.Auction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentAuction = auction
}

// UpdateAuction updates the current auction atomically
func (s *AuctionStore) UpdateAuction(updateFn func(*model.Auction) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentAuction == nil {
		return model.ErrNoActiveAuction
	}

	return updateFn(s.currentAuction)
}

// AddBid adds a bid to the current auction and returns the bid ID
func (s *AuctionStore) AddBid(bid *model.Bid) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentAuction == nil {
		return model.ErrNoActiveAuction
	}

	s.currentAuction.Bids = append(s.currentAuction.Bids, *bid)
	s.currentAuction.CurrentBid = bid.Amount
	s.currentAuction.CurrentWinner = &bid.UserID

	return nil
}

// GetNextBidID returns the next available bid ID
func (s *AuctionStore) GetNextBidID() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextBidID
	s.nextBidID++
	return id
}

// Subscribe creates a new subscription channel for auction events
func (s *AuctionStore) Subscribe(id string) chan *model.AuctionEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan *model.AuctionEvent, 10)
	s.subscribers[id] = ch
	return ch
}

// Unsubscribe removes a subscription
func (s *AuctionStore) Unsubscribe(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, exists := s.subscribers[id]; exists {
		close(ch)
		delete(s.subscribers, id)
	}
}

// Broadcast sends an event to all subscribers
func (s *AuctionStore) Broadcast(event *model.AuctionEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.subscribers {
		select {
		case ch <- event:
		default:
			// Skip slow consumers to prevent blocking
		}
	}
}

// GetSubscriberCount returns the number of active subscribers
func (s *AuctionStore) GetSubscriberCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}

// Clear removes the current auction (useful for testing)
func (s *AuctionStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentAuction = nil
}
