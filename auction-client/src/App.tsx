// src/App.tsx
import { ApolloProvider } from '@apollo/client';
import { client } from './apollo/client';
import AuctionDashboard from './components/AuctionDashboard';

function App() {
  return (
    <ApolloProvider client={client}>
      <AuctionDashboard />
    </ApolloProvider>
  );
}

export default App;