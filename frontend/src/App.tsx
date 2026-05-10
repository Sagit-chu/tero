import React, { useEffect, useState } from 'react';

function App() {
  const [status, setStatus] = useState('Loading...');

  useEffect(() => {
    fetch('/api/status')
      .then(r => r.json())
      .then(d => setStatus(d.status))
      .catch(() => setStatus('Error'));
  }, []);

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Flvx Monitor Dashboard</h1>
      <div className="bg-gray-100 p-4 rounded shadow">
        System Status: <span className="font-semibold">{status}</span>
      </div>
    </div>
  );
}

export default App;