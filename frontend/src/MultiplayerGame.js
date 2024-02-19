import React, { useEffect, useState } from 'react';
import { useParams, useLocation } from 'react-router-dom';

const MultiplayerGame = () => {
  const { sessionId } = useParams();
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const location = useLocation();
  const { playerName } = location.state; // Accessing player name passed in state

  console.log("Player name:", playerName);

  useEffect(() => {
    // Fetch questions based on sessionId
    // Update state with fetched questions
  }, [sessionId]);

  const handleAnswer = (selectedOption) => {
    // Submit answer
    // Possibly move to next question or handle game end
  };

  return (
    <div>
      <h1>Multiplayer Game</h1>
      {/* <p>Session ID: {sessionId}</p> */}
      <p>Question {currentQuestionIndex + 1} of {questions.length}</p>
      {/* Display question and options */}
    </div>
  );
};
export default MultiplayerGame;
