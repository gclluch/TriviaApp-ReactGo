// FinalScores.tsx
import React, { useEffect, useState } from 'react';
import { useParams, useLocation, useNavigate, Link } from 'react-router-dom';

interface LocationState {
  playerName: string;
  playerId: string;
  gameStarted: boolean;
}

interface Score {
  playerName: string; // Assuming this is the correct property based on your mapping
  score: number; // Assuming score is a numerical value
}

const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080';

const FinalScores: React.FC = () => {
  const { sessionId } = useParams<{ sessionId: string }>();
  const location = useLocation();
  const navigate = useNavigate();

  // Directly destructure and type assert the state to LocationState
  const { playerName } = (location.state || {
    playerName: 'Anonymous',
  }) as LocationState;

  const [scores, setScores] = useState<Score[]>([]);
  const [winners, setWinners] = useState<string[]>([]);
  const [highScore, setHighScore] = useState<number>(0);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchScores = async () => {
      try {
        const response = await fetch(`${API_BASE}/final-scores/${sessionId}`);
        if (!response.ok) {
          throw new Error('Failed to fetch scores');
        }
        const data = await response.json();
        setScores(data.scores);
        setWinners(data.winners);
        setHighScore(data.highScore);
      } catch (error: any) {
        // Type error as any if you're unsure of its structure
        setError(
          error.message ||
            'Failed to load final scores. Please try again later.'
        );
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
          <p key={index}>
            {score.playerName}: {score.score}
          </p>
        ))
      ) : (
        <p>No scores available.</p>
      )}
      <h2>Winners:</h2>
      {winners.length > 0 ? (
        winners.map((winner, index) => <p key={index}>{winner}</p>)
      ) : (
        <p>No winners.</p>
      )}
      <div>
        <Link to="/multiplayer">
          <button className="menu-button">Play Again</button>
        </Link>
      </div>
    </div>
  );
};

export default FinalScores;
