// LoadingIndicator.tsx
import React from 'react';

const LoadingIndicator: React.FC = () => (
  <div className="loading-container">
    <div aria-live="polite" className="loading">Loading...</div>
  </div>
);

export default LoadingIndicator;
