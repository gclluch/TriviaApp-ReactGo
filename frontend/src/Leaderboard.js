// Leaderboard.js

import React, { useState, useEffect } from 'react';
const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const Leaderboard = () => {
  const [leaderboard, setLeaderboard] = useState([]);

  useEffect(() => {
    const fetchLeaderboard = async () => {
        try {
            const response = await fetch('http://localhost:8080/leaderboard'); // Adjust URL as needed
            const data = await response.json();
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
                  <li key={index}>{playerID}: {score}</li>
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