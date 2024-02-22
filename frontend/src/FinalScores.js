import React, { useEffect, useState } from 'react';
import { useParams, useLocation, Link } from 'react-router-dom';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const FinalScores = () => {
  const { sessionId } = useParams();
  const location = useLocation();
  const { playerName } = location.state || { playerName: 'Anonymous' };
  const [scores, setScores] = useState([]);
  const [winners, setWinners] = useState([]);
  const [highScore, setHighScore] = useState(0);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchScores = async () => {
      try {
        const response = await fetch(`${API_BASE}/final-scores/${sessionId}`);
        if (!response.ok) throw new Error('Failed to fetch scores');
        const data = await response.json();
        setScores(data.scores);
        setWinners(data.winners);
        setHighScore(data.highScore);
      } catch (error) {
        setError('Failed to load final scores. Please try again later.');
      }
    };

    fetchScores();
  }, [sessionId]);

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div>
      <h3>Final Scores</h3>
      <p>High Score: {highScore}</p>
      <p>You: {playerName}</p>
      {scores.length > 0 ? (
        scores.map((score, index) => (
          <p key={index}>{score.playerName}: {score.score}</p>
        ))
      ) : (
        <p>No scores available.</p>
      )}
      <h2>Winners:</h2>
      {winners.length > 0 ? (
        winners.map((winner, index) => (
          <p key={index}>{winner}</p>
        ))
      ) : (
        <p>No winners.</p>
      )}
      <div>
        <Link to="/multiplayer"><button className="menu-button">Play Again</button></Link>
      </div>
    </div>
  );
};

export default FinalScores;
