// Leaderboard.tsx

import React, { useState, useEffect } from 'react';

// Define TypeScript interface for the leaderboard entries
interface LeaderboardEntry {
  [playerID: string]: string; // Assuming score is stored as a string. Adjust if it's a different type.
}

const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080';

const Leaderboard: React.FC = () => {
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry>({});

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        const response = await fetch(`${API_BASE}/leaderboard`); // Use API_BASE for consistency
        const data: LeaderboardEntry = await response.json(); // Type assertion
        setLeaderboard(data);
      } catch (error) {
        console.error('Failed to fetch leaderboard:', error);
      }
    };

    fetchLeaderboard();
  }, []);

  // Function to render the leaderboard if it has been set
  const renderLeaderboard = () => {
    if (!leaderboard || Object.keys(leaderboard).length === 0) {
      return <p>No leaderboard data available.</p>;
    }

    return (
      <ul>
        {Object.entries(leaderboard).map(([playerID, score], index) => (
          <li key={index}>
            {playerID}: {score}
          </li>
        ))}
      </ul>
    );
  };

  return (
    <div>
      <h2>Leaderboard</h2>
      {renderLeaderboard()}
    </div>
  );
};

export default Leaderboard;
