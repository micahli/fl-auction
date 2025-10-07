import React, { useState, useEffect, useRef } from 'react';
import { useMutation, useSubscription, useQuery } from '@apollo/client';
import {
  CREATE_AUCTION,
  PLACE_BID,
  GET_CURRENT_AUCTION,
  AUCTION_EVENTS_SUBSCRIPTION,
} from '../graphql/operations';
import { Timer, DollarSign, User, AlertCircle, Play, TrendingUp, Settings } from 'lucide-react';

const AuctionDashboard: React.FC = () => {
  const [userId] = useState(`User${Math.floor(Math.random() * 1000)}`);
  const [bidAmount, setBidAmount] = useState('');
  const [bidError, setBidError] = useState('');
  const [bidSuccess, setBidSuccess] = useState(false);
  const [auctionExtended, setAuctionExtended] = useState(false);
  const [auctionData, setAuctionData] = useState<any>(null);
  const [localTimeRemaining, setLocalTimeRemaining] = useState<number>(0);
  const previousTimeRef = useRef<number>(0);
  
  // Auction creation form state
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [startingBid, setStartingBid] = useState('100');
  const [duration, setDuration] = useState('30');
  const [extendedBidding, setExtendedBidding] = useState(true);

  // GraphQL operations
  const { data: queryData, refetch } = useQuery(GET_CURRENT_AUCTION);
  
  const [createAuction, { loading: creatingAuction }] = useMutation(CREATE_AUCTION, {
    onCompleted: (data) => {
      setAuctionData(data.createAuction);
      setLocalTimeRemaining(data.createAuction.timeRemaining);
      setShowCreateForm(false);
      // Reset form
      setStartingBid('100');
      setDuration('30');
      setExtendedBidding(true);
    },
    onError: (error) => {
      setBidError(error.message);
    },
  });

  const [placeBid, { loading: placingBid }] = useMutation(PLACE_BID, {
    onCompleted: () => {
      setBidSuccess(true);
      setTimeout(() => setBidSuccess(false), 3000);
      refetch();
    },
    onError: (error) => {
      setBidError(error.message);
    },
  });

  const { data: subData } = useSubscription(AUCTION_EVENTS_SUBSCRIPTION, {
    onData: ({ data }) => {
      const event = data.data?.auctionEvents;
      if (event?.auction) {
        const newTimeRemaining = event.auction.timeRemaining;
        const oldTimeRemaining = previousTimeRef.current;
        
        console.log(`📊 Subscription update - Type: ${event.type}, Old: ${oldTimeRemaining}s, New: ${newTimeRemaining}s`);
        
        // Update auction data
        setAuctionData(event.auction);
        
        // Check if time was extended (time increased instead of decreasing)
        if (event.type === 'BID_PLACED' && newTimeRemaining > oldTimeRemaining) {
          console.log(`🔔 AUCTION EXTENDED! ${oldTimeRemaining}s → ${newTimeRemaining}s`);
          setAuctionExtended(true);
          setTimeout(() => setAuctionExtended(false), 3000);
        }
        
        // Update the local timer with the new value from backend
        setLocalTimeRemaining(newTimeRemaining);
        previousTimeRef.current = newTimeRemaining;
      }
    },
  });

  // Initialize auction data from query
  useEffect(() => {
    if (queryData?.currentAuction) {
      setAuctionData(queryData.currentAuction);
      setLocalTimeRemaining(queryData.currentAuction.timeRemaining);
      previousTimeRef.current = queryData.currentAuction.timeRemaining;
      setBidAmount(queryData.currentAuction.nextBid.toString());
    }
  }, [queryData]);

  // Update bid amount when auction changes
  useEffect(() => {
    if (auctionData?.nextBid) {
      setBidAmount(auctionData.nextBid.toString());
    }
  }, [auctionData?.nextBid]);

  // Local countdown timer - updates every second
  useEffect(() => {
    if (!auctionData || auctionData.status !== 'ACTIVE') {
      return;
    }

    const interval = setInterval(() => {
      setLocalTimeRemaining((prev) => {
        const newValue = prev - 1;
        if (newValue <= 0) {
          clearInterval(interval);
          refetch();
          return 0;
        }
        previousTimeRef.current = newValue;
        return newValue;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [auctionData?.status, auctionData?.id, refetch]);

  const handleCreateAuction = async () => {
    setBidError('');
    
    const startBid = parseFloat(startingBid);
    const dur = parseInt(duration);

    if (isNaN(startBid) || startBid <= 0) {
      setBidError('Starting bid must be a positive number');
      return;
    }

    if (isNaN(dur) || dur < 10 || dur > 3600) {
      setBidError('Duration must be between 10 and 3600 seconds');
      return;
    }

    try {
      await createAuction({
        variables: {
          startingBid: startBid,
          duration: dur,
          extendedBidding: extendedBidding,
        },
      });
    } catch (err) {
      // Error handled in onError
    }
  };

  const handlePlaceBid = async () => {
    setBidError('');
    setBidSuccess(false);

    const amount = parseFloat(bidAmount);
    if (isNaN(amount)) {
      setBidError('Please enter a valid bid amount');
      return;
    }

    try {
      await placeBid({
        variables: {
          userId,
          amount,
        },
      });
    } catch (err) {
      // Error handled in onError
    }
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  // Show create form when: no auction OR user clicked to show form
  const shouldShowCreateForm = !auctionData || showCreateForm;

  // RENDER: Create Auction Form
  if (shouldShowCreateForm) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex flex-col items-center justify-center gap-6 p-6">
        <h1 className="text-4xl font-bold text-white mb-2">Live Auction System</h1>
        <p className="text-gray-400 mb-6">No active auction</p>

        <div className="bg-gray-800 rounded-2xl shadow-2xl p-8 border border-gray-700 max-w-md w-full">
          <div className="flex items-center gap-3 mb-6">
            <Settings className="text-blue-500" size={24} />
            <h2 className="text-2xl font-bold text-white">Create Auction</h2>
          </div>

          <div className="space-y-6">
            {/* Starting Bid */}
            <div>
              <label className="block text-gray-300 text-sm font-semibold mb-2">
                Starting Bid ($)
              </label>
              <div className="relative">
                <DollarSign
                  className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                  size={20}
                />
                <input
                  type="number"
                  value={startingBid}
                  onChange={(e) => setStartingBid(e.target.value)}
                  className="w-full bg-gray-900 text-white pl-10 pr-4 py-3 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 transition"
                  placeholder="100"
                  step="1"
                  min="1"
                />
              </div>
              <p className="text-gray-500 text-xs mt-1">Minimum: $1</p>
            </div>

            {/* Duration */}
            <div>
              <label className="block text-gray-300 text-sm font-semibold mb-2">
                Duration (seconds)
              </label>
              <div className="relative">
                <Timer
                  className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                  size={20}
                />
                <input
                  type="number"
                  value={duration}
                  onChange={(e) => setDuration(e.target.value)}
                  className="w-full bg-gray-900 text-white pl-10 pr-4 py-3 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 transition"
                  placeholder="30"
                  step="5"
                  min="10"
                  max="3600"
                />
              </div>
              <p className="text-gray-500 text-xs mt-1">Range: 10-3600 seconds</p>
            </div>

            {/* Extended Bidding Toggle */}
            <div className="bg-gray-900/50 rounded-lg p-4 border border-gray-700">
              <div className="flex items-center justify-between mb-2">
                <label className="text-gray-300 text-sm font-semibold">
                  Extended Bidding
                </label>
                <button
                  onClick={() => setExtendedBidding(!extendedBidding)}
                  className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                    extendedBidding ? 'bg-blue-600' : 'bg-gray-600'
                  }`}
                >
                  <span
                    className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                      extendedBidding ? 'translate-x-6' : 'translate-x-1'
                    }`}
                  />
                </button>
              </div>
              <p className="text-gray-400 text-xs">
                {extendedBidding
                  ? '✓ Auction extends by 10 seconds if bid placed in final 10 seconds'
                  : '✗ Auction ends at scheduled time regardless of bids'}
              </p>
            </div>

            {/* Error Message */}
            {bidError && (
              <div className="bg-red-500/10 border border-red-500 rounded-lg p-3 flex items-center gap-2 text-red-400 text-sm">
                <AlertCircle size={16} />
                {bidError}
              </div>
            )}

            {/* Action Buttons */}
            <div className="flex gap-3">
              {auctionData && (
                <button
                  onClick={() => {
                    setShowCreateForm(false);
                    setBidError('');
                  }}
                  className="flex-1 bg-gray-700 hover:bg-gray-600 text-white px-6 py-3 rounded-lg font-semibold transition"
                >
                  Cancel
                </button>
              )}
              <button
                onClick={handleCreateAuction}
                disabled={creatingAuction}
                className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-800 text-white px-6 py-3 rounded-lg font-semibold transition"
              >
                {creatingAuction ? 'Creating...' : 'Create Auction'}
              </button>
            </div>
          </div>
        </div>

        <div className="mt-4 bg-gray-800/50 rounded-lg p-4 max-w-md border border-gray-700">
          <div className="flex items-start gap-3 text-gray-400 text-sm">
            <AlertCircle size={18} className="flex-shrink-0 mt-0.5" />
            <p>
              Your User ID: <span className="text-white font-semibold">{userId}</span>
            </p>
          </div>
        </div>
      </div>
    );
  }

  // RENDER: Active or Ended Auction
  const isActive = auctionData.status === 'ACTIVE';
  const isEnding = localTimeRemaining <= 10;

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">Live Auction</h1>
          <p className="text-gray-400">
            Your ID: <span className="text-blue-400 font-semibold">{userId}</span>
          </p>
        </div>

        {/* Main Auction Card */}
        <div className="bg-gray-800 rounded-2xl shadow-2xl overflow-hidden border border-gray-700">
          {/* Status Banner */}
          <div
            className={`py-3 px-6 ${
              isActive
                ? isEnding
                  ? 'bg-orange-600'
                  : 'bg-green-600'
                : 'bg-gray-700'
            }`}
          >
            <div className="flex items-center justify-between">
              <span className="text-white font-semibold flex items-center gap-2">
                <div
                  className={`w-2 h-2 rounded-full ${
                    isActive ? 'bg-white animate-pulse' : 'bg-gray-400'
                  }`}
                ></div>
                {isActive
                  ? isEnding
                    ? 'Ending Soon!'
                    : 'Auction Active'
                  : 'Auction Ended'}
                {auctionData.extendedBidding && isActive && (
                  <span className="ml-2 text-xs bg-white/20 px-2 py-1 rounded">
                    Extended Bidding ON
                  </span>
                )}
              </span>
              <div className="flex items-center gap-2 text-white">
                <Timer size={18} />
                <span
                  className={`text-lg font-mono font-bold ${
                    isEnding ? 'animate-pulse' : ''
                  }`}
                >
                  {formatTime(localTimeRemaining)}
                </span>
              </div>
            </div>
          </div>

          {/* Current Bid Section */}
          <div className="p-8">
            <div className="text-center mb-8">
              <div className="text-gray-400 text-sm mb-2">CURRENT BID</div>
              <div className="text-6xl font-bold text-white flex items-center justify-center gap-2">
                <DollarSign size={48} className="text-green-500" />
                {auctionData.currentBid.toFixed(2)}
              </div>

              {auctionData.currentWinner && (
                <div className="mt-4 flex items-center justify-center gap-2 text-gray-300">
                  <User size={18} className="text-blue-500" />
                  <span>
                    Current Winner:{' '}
                    <span className="font-semibold text-white">
                      {auctionData.currentWinner}
                    </span>
                  </span>
                </div>
              )}

              {!auctionData.currentWinner && (
                <div className="mt-4 text-gray-400">
                  No bids yet - Starting bid: ${auctionData.startingBid}
                </div>
              )}
            </div>

            {/* Bid Form - Only show when auction is active */}
            {isActive && (
              <div className="max-w-md mx-auto">
                <div className="bg-gray-700/50 rounded-xl p-6 border border-gray-600">
                  <label className="block text-gray-300 text-sm font-semibold mb-3">
                    Place Your Bid
                  </label>

                  <div className="flex gap-3 mb-4">
                    <div className="flex-1 relative">
                      <DollarSign
                        className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                        size={20}
                      />
                      <input
                        type="number"
                        value={bidAmount}
                        onChange={(e) => setBidAmount(e.target.value)}
                        className="w-full bg-gray-900 text-white pl-10 pr-4 py-3 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 transition"
                        placeholder="Enter bid"
                        step="0.01"
                        disabled={placingBid}
                      />
                    </div>
                    <button
                      onClick={handlePlaceBid}
                      disabled={placingBid}
                      className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-800 text-white px-6 py-3 rounded-lg font-semibold transition transform hover:scale-105 active:scale-95 shadow-lg"
                    >
                      {placingBid ? 'Bidding...' : 'Bid'}
                    </button>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <span className="text-gray-400">Next minimum bid:</span>
                    <span className="text-green-400 font-semibold flex items-center gap-1">
                      <TrendingUp size={14} />
                      ${auctionData.nextBid.toFixed(2)}
                    </span>
                  </div>

                  {bidError && (
                    <div className="mt-4 bg-red-500/10 border border-red-500 rounded-lg p-3 flex items-center gap-2 text-red-400 text-sm">
                      <AlertCircle size={16} />
                      {bidError}
                    </div>
                  )}

                  {bidSuccess && (
                    <div className="mt-4 bg-green-500/10 border border-green-500 rounded-lg p-3 text-green-400 text-sm font-semibold text-center">
                      ✓ Bid placed successfully!
                    </div>
                  )}

                  {auctionExtended && (
                    <div className="mt-4 bg-orange-500/10 border border-orange-500 rounded-lg p-3 text-orange-400 text-sm font-semibold text-center animate-pulse">
                      ⏰ Auction extended by 10 seconds!
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Auction Ended */}
            {!isActive && (
              <div className="text-center">
                <div className="bg-gray-700/50 rounded-xl p-8 border border-gray-600">
                  <h3 className="text-2xl font-bold text-white mb-4">
                    Auction Ended
                  </h3>
                  {auctionData.currentWinner ? (
                    <div>
                      <p className="text-gray-300 mb-2">Winner:</p>
                      <div className="text-3xl font-bold text-green-500 flex items-center justify-center gap-2 mb-4">
                        <User size={32} />
                        {auctionData.currentWinner}
                      </div>
                      <p className="text-gray-400">
                        Winning bid:{' '}
                        <span className="text-white font-semibold">
                          ${auctionData.currentBid.toFixed(2)}
                        </span>
                      </p>
                    </div>
                  ) : (
                    <p className="text-gray-400">No bids were placed</p>
                  )}
                  
                  <button
                    onClick={() => setShowCreateForm(true)}
                    className="mt-6 bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-semibold transition"
                  >
                    Start New Auction
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Info Card */}
        <div className="mt-6 bg-gray-800/50 rounded-lg p-4 border border-gray-700">
          <div className="flex items-start gap-3 text-gray-400 text-sm">
            <AlertCircle size={18} className="flex-shrink-0 mt-0.5" />
            <p>
              Real-time updates powered by{' '}
              <span className="text-blue-400 font-semibold">
                GraphQL Subscriptions
              </span>
              . Open multiple browser tabs to see live synchronization!
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuctionDashboard;