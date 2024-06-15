import React, { Suspense } from 'react';
import RouteList from './routes/index';

function App() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <RouteList />
    </Suspense>
  );
}

export default App;
