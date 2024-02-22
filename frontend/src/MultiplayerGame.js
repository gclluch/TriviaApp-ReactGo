import React, { useEffect, useState } from 'react';
import { useParams, useLocation, useNavigate } from 'react-router-dom';
import QuestionDisplay from './QuestionDisplay';
import ScoreDisplay from './ScoreDisplay';
import { useWebSocket } from './WebSocketContext';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
  const { sessionId } = useParams();
  const navigate = useNavigate();
  const { playerName, playerId, gameStarted } = useLocation().state || {};
  const [failedToJoin, setFailedToJoin] = useState(!gameStarted);
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);
  const [highScore, setHighScore] = useState(0);
  const { webSocket, isConnected } = useWebSocket();
  const [hasFinished, setHasFinished] = useState(false);

  useEffect(() => {
    const fetchQuestions = async () => {
      try {
        const response = await fetch(`${API_BASE}/questions/${sessionId}`);
        const data = await response.json();
        setQuestions(data.questions);
      } catch (error) {
        console.error("Failed to fetch questions:", error);
      }
    };

    if (sessionId && gameStarted) {
      fetchQuestions();
    }
  }, [sessionId, gameStarted]);

  useEffect(() => {
    const handleMessage = (event) => {
      const data = JSON.parse(event.data);
      switch (data.type) {
        case 'highScore':
          setHighScore(data.score);
          break;
        case 'sessionComplete':
          navigate(`/final-scores/${sessionId}`, { state: { playerName, playerId } });
          break;
        default:
          console.log("Unhandled message type:", data.type);
      }
    };

    if (webSocket && isConnected) {
      webSocket.addEventListener('message', handleMessage);
    }

    return () => webSocket?.removeEventListener('message', handleMessage);
  }, [webSocket, isConnected, sessionId, navigate, playerName, playerId]);

  const submitAnswer = async (index) => {
    const currentQuestion = questions[currentQuestionIndex];
    try {
      const response = await fetch(`${API_BASE}/answer`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          sessionId,
          playerId,
          questionId: currentQuestion.id,
          answer: index,
        }),
      });
      const data = await response.json();

      if (data.correct) {
        setScore(prevScore => prevScore + 10); // Assuming 10 points per correct answer
      }

      const nextIndex = currentQuestionIndex + 1;
      if (nextIndex < questions.length) {
        setCurrentQuestionIndex(nextIndex);
      } else {
        markPlayerFinished();
      }
    } catch (error) {
      console.error("Failed to submit answer:", error);
    }
  };

  const markPlayerFinished = async () => {
    try {
      await fetch(`${API_BASE}/player/finished`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sessionId, playerId }),
      });
      setHasFinished(true);
    } catch (error) {
      console.error("Error marking player as finished:", error);
    }
  };

  if (failedToJoin) {
    return <div><h3>Failed to Join Game in Time</h3></div>;
  }

  return (
    <div>
      <h1>Multiplayer Game</h1>
      {!hasFinished ? (
        <>
          <p>Question {currentQuestionIndex + 1} of {questions.length}</p>
          <QuestionDisplay
            question={questions[currentQuestionIndex]?.questionText}
            options={questions[currentQuestionIndex]?.options}
            onAnswer={(index) => submitAnswer(index)} // Pass index to submitAnswer directly
          />
        </>
      ) : (
        <div><h2>You've completed all questions. Waiting for other players to finish...</h2></div>
      )}
      <ScoreDisplay score={score} highScore={highScore} />
    </div>
  );
};

export default MultiplayerGame;
