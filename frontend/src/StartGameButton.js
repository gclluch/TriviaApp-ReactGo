// Code for the start game button

import React from 'react';
import PropTypes from 'prop-types';

const StartGameButton = ({onStart}) => (
    <button onClick={onStart}>Start Game</button>
);

StartGameButton.propTypes = {
    onStart: PropTypes.func.isRequired,
};

export default StartGameButton;