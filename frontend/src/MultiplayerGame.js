import React, { useEffect, useState } from 'react';
import { useParams, useLocation, useNavigate } from 'react-router-dom';
import QuestionDisplay from './QuestionDisplay';
import ScoreDisplay from './ScoreDisplay';
const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
  const { sessionId } = useParams();
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const location = useLocation();
  const navigate = useNavigate();

  const [score, setScore] = useState(0);
  const [highSCore, setHighScore] = useState(0);
  const { playerName, gameStarted } = location.state || { gameStarted: true }; // Default to gameStarted if state is undefined

  const [failedToJoin, setFailedToJoin] = useState(false);


  console.log("Player name:", playerName);

  useEffect(() => {
    if (!gameStarted) {
      // If the gameStarted flag is false, it means the player failed to join in time
      setFailedToJoin(true);
      // Optional: navigate back or to another page after a delay
      setTimeout(() => {
        navigate('/');
      }, 5000); // Redirect to home or another appropriate page after 5 seconds
    }
  }, [gameStarted, navigate]);

  useEffect(() => {
    // Fetch questions based on sessionId
    // Update state with fetched questions
    fetchQuestions(sessionId)
  }, [sessionId]);

  if (failedToJoin) {
    return <div><h3>Failed to Join Game in Time</h3></div>;
  }

  const fetchQuestions = async (session) => {
    try {
      const response = await fetch(`${API_BASE}/questions`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sessionId: session, // need to provide the sessionId
        }),
      });

      const data = await response.json();
      setQuestions(data);
    } catch (error) {
      // setError("Failed to fetch questions.");
    }
  };

  const handleAnswer = (selectedOption) => {
    // Submit answer
    // Possibly move to next question or handle game end
  };

  const submitAnswer = (index) => async () => {};

  return (
    <div>
      <h1>Multiplayer Game</h1>
      {/* <p>Session ID: {sessionId}</p> */}
      <p>Question {currentQuestionIndex + 1} of {questions.length}</p>
      <QuestionDisplay
          question={questions[currentQuestionIndex]?.questionText}
          options={questions[currentQuestionIndex]?.options}
          onAnswer={submitAnswer}
        />
      {/* Display question and options */}
    </div>
  );
};
export default MultiplayerGame;
