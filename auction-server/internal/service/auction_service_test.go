package service

import (
	"context"
	"testing"
	"time"

	"github.com/micahli/fl-auction/auction-server/internal/model"
	"github.com/micahli/fl-auction/auction-server/internal/store"
)

func TestCreateAuction(t *testing.T) {
	st := store.NewAuctionStore()
	svc := NewAuctionService(st)

	auction, err := svc.CreateAuction(context.Background(), 100.0, 30, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if auction.StartingBid != 100.0 {
		t.Errorf("expected starting bid 100.0, got %f", auction.StartingBid)
	}

	if auction.Status != model.AuctionStatusActive {
		t.Errorf("expected status ACTIVE, got %s", auction.Status)
	}
}

func TestCreateAuction_AlreadyActive(t *testing.T) {
	st := store.NewAuctionStore()
	svc := NewAuctionService(st)

	_, err := svc.CreateAuction(context.Background(), 100.0, 30, false)
	if err != nil {
		t.Fatalf("first auction creation failed: %v", err)
	}

	_, err = svc.CreateAuction(context.Background(), 100.0, 30, false)
	if err != model.ErrAuctionAlreadyActive {
		t.Errorf("expected ErrAuctionAlreadyActive, got %v", err)
	}
}

func TestPlaceBid_Success(t *testing.T) {
	st := store.NewAuctionStore()
	svc := NewAuctionService(st)

	_, err := svc.CreateAuction(context.Background(), 100.0, 30, false)
	if err != nil {
		t.Fatalf("auction creation failed: %v", err)
	}

	bid, err := svc.PlaceBid(context.Background(), "user1", 150.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bid.Amount != 150.0 {
		t.Errorf("expected bid amount 150.0, got %f", bid.Amount)
	}

	if bid.UserID != "user1" {
		t.Errorf("expected user1, got %s", bid.UserID)
	}
}

func TestPlaceBid_TooLow(t *testing.T) {
	st := store.NewAuctionStore()
	svc := NewAuctionService(st)

	_, err := svc.CreateAuction(context.Background(), 100.0, 30, false)
	if err != nil {
		t.Fatalf("auction creation failed: %v", err)
	}

	_, err = svc.PlaceBid(context.Background(), "user1", 100.0)
	if err == nil {
		t.Error("expected bid too low error")
	}
}

func TestPlaceBid_ExtendedBidding(t *testing.T) {
	st := store.NewAuctionStore()
	svc := NewAuctionService(st)

	auction, err := svc.CreateAuction(context.Background(), 100.0, 5, true)
	if err != nil {
		t.Fatalf("auction creation failed: %v", err)
	}

	originalEndTime := auction.EndTime

	// Wait until < 10 seconds remaining
	time.Sleep(1 * time.Second)

	_, err = svc.PlaceBid(context.Background(), "user1", 150.0)
	if err != nil {
		t.Fatalf("bid placement failed: %v", err)
	}

	currentAuction := st.GetCurrentAuction()
	if !currentAuction.EndTime.After(originalEndTime) {
		t.Error("expected auction to be extended")
	}
}

func TestAuctionExpiry(t *testing.T) {
	st := store.NewAuctionStore()
	svc := NewAuctionService(st)

	_, err := svc.CreateAuction(context.Background(), 100.0, 2, false)
	if err != nil {
		t.Fatalf("auction creation failed: %v", err)
	}

	// Wait for auction to expire
	time.Sleep(3 * time.Second)

	auction := st.GetCurrentAuction()
	if auction.Status != model.AuctionStatusEnded {
		t.Errorf("expected status ENDED, got %s", auction.Status)
	}

	// Try to place bid on ended auction
	_, err = svc.PlaceBid(context.Background(), "user1", 150.0)
	if err != model.ErrNoActiveAuction {
		t.Errorf("expected ErrNoActiveAuction, got %v", err)
	}
}
