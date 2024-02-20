// ScoreDisplay.js

import React from 'react';
import PropTypes from 'prop-types';

const ScoreDisplay = ({ score, highScore }) => (
    <div>
        <p className="score">Your Score: {score}</p>
        {highScore !== undefined && <p className="score">Top Score: {highScore}</p>}
    </div>
);

ScoreDisplay.propTypes = {
    score: PropTypes.number.isRequired,
    highScore: PropTypes.number, // highScore is optional
};

ScoreDisplay.defaultProps = {
    highScore: undefined, // Not required for single player, default to undefined
};

export default ScoreDisplay;
