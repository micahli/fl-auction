# fl-auction
take home task for fl's interview

// Example GraphQL Queries/Mutations/Subscriptions:
/*

# Create an auction
mutation {
  createAuction(startingBid: 100.0, duration: 30, extendedBidding: true) {
    id
    currentBid
    status
    endTime
  }
}

# Place a bid
mutation {
  placeBid(userId: "user123", amount: 150.0) {
    id
    amount
    timestamp
  }
}

# Query current auction
query {
  currentAuction {
    id
    currentBid
    currentWinner
    nextBid
    timeRemaining
    status
  }
}

# Subscribe to auction events
subscription {
  auctionEvents {
    type
    auction {
      id
      currentBid
      currentWinner
      nextBid
      timeRemaining
      status
    }
    bid {
      userId
      amount
    }
  }
}

*/