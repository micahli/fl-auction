import { gql } from '@apollo/client';

// Mutation: Create a new auction
export const CREATE_AUCTION = gql`
  mutation CreateAuction(
    $startingBid: Float!
    $duration: Int
    $extendedBidding: Boolean
  ) {
    createAuction(
      startingBid: $startingBid
      duration: $duration
      extendedBidding: $extendedBidding
    ) {
      id
      startingBid
      currentBid
      currentWinner
      duration
      extendedBidding
      status
      nextBid
      timeRemaining
    }
  }
`;

// Mutation: Place a bid
export const PLACE_BID = gql`
  mutation PlaceBid($userId: String!, $amount: Float!) {
    placeBid(userId: $userId, amount: $amount) {
      id
      userId
      amount
      timestamp
    }
  }
`;

// Query: Get current auction
export const GET_CURRENT_AUCTION = gql`
  query GetCurrentAuction {
    currentAuction {
      id
      startingBid
      currentBid
      currentWinner
      duration
      extendedBidding
      status
      nextBid
      timeRemaining
    }
  }
`;

// Subscription: Real-time auction events
export const AUCTION_EVENTS_SUBSCRIPTION = gql`
  subscription AuctionEvents {
    auctionEvents {
      type
      auction {
        id
        currentBid
        currentWinner
        status
        nextBid
        timeRemaining
      }
      bid {
        id
        userId
        amount
        timestamp
      }
      error
    }
  }
`;