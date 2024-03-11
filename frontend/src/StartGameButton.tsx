// StartGameButton.tsx
import React from 'react';

// Define an interface for the component's props
interface StartGameButtonProps {
  onStart: () => void; // Specifies that onStart is a function that takes no parameters and returns void
}

const StartGameButton: React.FC<StartGameButtonProps> = ({ onStart }) => (
  <button onClick={onStart}>Start Game</button>
);

export default StartGameButton;
