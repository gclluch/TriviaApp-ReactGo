import React from 'react';
import PropTypes from 'prop-types';

const ScoreDisplay = ({ score, highScore }) => (
    <div className="score-display">
        <p className="score">Your Score: {score}</p>
        {highScore !== null && <p className="score">Top Score: {highScore}</p>}
    </div>
);

ScoreDisplay.propTypes = {
  score: PropTypes.number.isRequired,
  highScore: PropTypes.number, // highScore is optional
};

ScoreDisplay.defaultProps = {
  highScore: null, // Default to null for clarity when highScore is not applicable
};

export default ScoreDisplay;
