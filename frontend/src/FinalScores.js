// FinalScores.js
import React, { useEffect, useState } from 'react';
import { useParams, useLocation } from 'react-router-dom';
const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const FinalScores = () => {
  const { sessionId } = useParams();
  const location = useLocation();
  const { playerName, playerId } = location.state || {};
  const [scores, setScores] = useState([]);
  const [winners, setWinners] = useState([]);
  const [highScore, setHighScore] = useState(0);

  useEffect(() => {
    const fetchScores = async () => {
      const response = await fetch(`${API_BASE}/final-scores/${sessionId}`);
      const data = await response.json();
      setScores(data.scores);
      setWinners(data.winners);
      setHighScore(data.highScore);
    };

    fetchScores();
  }, [sessionId]);

  return (
    <div>
      <h1>Final Scores</h1>
      <p>High Score: {highScore}</p>
      <p>You: {playerName}</p>
      {scores.map((score, index) => (
        <p key={index}>{score.playerName}: {score.score}</p>
      ))}
      <h2>Winners:</h2>
      {winners.map((winner, index) => (
        <p key={index}>{winner}</p>
      ))}
    </div>
  );
};

export default FinalScores;