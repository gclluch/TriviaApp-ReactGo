// ScoreDisplay.js

import React from 'react';
import PropTypes from 'prop-types';

const ScoreDisplay = ({ score }) => (
    <div>
        <p className="score">Score: {score}</p>
    </div>
);

ScoreDisplay.propTypes = {
    score: PropTypes.number.isRequired,
};

export default ScoreDisplay;
