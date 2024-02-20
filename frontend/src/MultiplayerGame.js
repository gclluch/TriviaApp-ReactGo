import React, { useEffect, useState } from 'react';
import { useParams, useLocation, useNavigate } from 'react-router-dom';
import QuestionDisplay from './QuestionDisplay';
import ScoreDisplay from './ScoreDisplay';
import { useWebSocket } from './WebSocketContext'; // Ensure this is correctly implemented

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
  const { sessionId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const { playerName, playerId, gameStarted } = location.state || {};
  const [failedToJoin, setFailedToJoin] = useState(false);

  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);
  const [highScore, setHighScore] = useState(0); // State to store the high score
  const { webSocket, isConnected } = useWebSocket(); // Destructuring to get webSocket and isConnected
  const [hasFinished, setHasFinished] = useState(false);



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
    const fetchQuestions = async () => {
      try {
        const response = await fetch(`${API_BASE}/questions`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ sessionId }),
        });
        const data = await response.json();
        setQuestions(data);
      } catch (error) {
        console.error("Failed to fetch questions:", error);
      }
    };

    if (sessionId) {
      fetchQuestions();
    }
  }, [sessionId]);

  useEffect(() => {
    if (webSocket && isConnected) {
      const handleMessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'highScore') {
          setHighScore(data.score); // Update high score on receiving a highScore message
        }
      };

      webSocket.addEventListener('message', handleMessage);

      return () => {
        webSocket.removeEventListener('message', handleMessage);
      };
    }
  }, [webSocket, isConnected]);

  if (failedToJoin) {
    return <div><h3>Failed to Join Game in Time</h3></div>;
  }



  const submitAnswer = (index) => async () => {
    const currentQuestion = questions[currentQuestionIndex];
    try {
      const res = await fetch(`${API_BASE}/answer`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          sessionId,
          playerId,
          questionId: currentQuestion.id,
          answer: index,
        }),
      });
      const data = await res.json();

      if (data.correct) {
        setScore(data.currentScore); // Update score if points were added
    }
    if (currentQuestionIndex < questions.length - 1) {
        setCurrentQuestionIndex(currentQuestionIndex + 1);
    } else {
        // Handle end game scenario for the player
        setHasFinished(true);
    }
} catch (error) {
    console.error("Failed to submit answer:", error);
}
  };


  return (
    <div>
      <h1>Multiplayer Game</h1>
      <p>Question {currentQuestionIndex + 1} of {questions.length}</p>
      {questions.length > 0 && (
        <QuestionDisplay
          question={questions[currentQuestionIndex].questionText}
          options={questions[currentQuestionIndex].options}
          onAnswer={submitAnswer}
        />
      )}
      <ScoreDisplay score={score} highScore={highScore} />
    </div>
  );
};
export default MultiplayerGame;
