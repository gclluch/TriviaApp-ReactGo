// ScoreDisplay.tsx
import React from 'react';

// Define an interface for the component props
interface ScoreDisplayProps {
  score: number;
  highScore?: number; // highScore is optional
}

const ScoreDisplay: React.FC<ScoreDisplayProps> = ({
  score,
  highScore = null,
}) => (
  <div className="score-display">
    <p className="score">Your Score: {score}</p>
    {highScore !== null && <p className="score">Top Score: {highScore}</p>}
  </div>
);

export default ScoreDisplay;
