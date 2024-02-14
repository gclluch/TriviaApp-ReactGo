import React, { useState } from "react";
import QuestionDisplay from './QuestionDisplay';
import ScoreDisplay from './ScoreDisplay';
import LoadingIndicator from './LoadingIndicator';
import ErrorMessage from './ErrorMessage';
import StartGameButton from './StartGameButton';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const SinglePlayerGame = () => {
  const [gameSession, setGameSession] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const startGame = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      const data = await response.json();
      setGameSession(data.sessionId);
      await fetchQuestions();
    } catch (error) {
      setError("Failed to start game.");
      setLoading(false);
    }
  };

  const fetchQuestions = async () => {
    try {
      const response = await fetch(`${API_BASE}/questions`);
      const data = await response.json();
      setQuestions(data);
      setLoading(false); // Stop loading once questions are fetched
    } catch (error) {
      setError("Failed to fetch questions.");
      setLoading(false);
    }
  };

  const submitAnswer = (index) => async () => {
    // setLoading(true);
    const currentQuestion = questions[currentQuestionIndex];
    try {
      const res = await fetch(`${API_BASE}/answer`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sessionId: gameSession,
          questionId: currentQuestion.id,
          answer: index,
        }),
      });
      const data = await res.json();
      if (data.correct) {
        setScore(data.currentScore);
      }
      if (currentQuestionIndex < questions.length - 1) {
        setCurrentQuestionIndex(currentQuestionIndex + 1);
      } else {
        endGame();
      }
    } catch (err) {
      setError("Failed to submit answer.");
    }
    // setLoading(false);
  };

  const endGame = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/game/end`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sessionId: gameSession, // need to provide the sessionId
        }),
      });
      const data = await res.json();
      alert(`Game over! Your score: ${data.finalScore}`); // Use the finalScore from the response
      setGameSession(null);
      setQuestions([]);
      setCurrentQuestionIndex(0);
      setScore(0);
    } catch (err) {
      setError("Failed to end game.");
    }
    setLoading(false);
  };

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage message={error} />;

  return (
    <div>
      <h1>Single Player</h1>
      {!gameSession ? (
        <StartGameButton onStart={startGame} />
      ) : currentQuestionIndex < questions.length ? (
        <QuestionDisplay
          question={questions[currentQuestionIndex]?.questionText}
          options={questions[currentQuestionIndex]?.options}
          onAnswer={submitAnswer}
        />
      ) : (
        <div>
          Game Over. Your score: {score}.
          <button onClick={() => setGameSession(null)}>Restart</button>
        </div>
      )}
      {gameSession && <ScoreDisplay score={score}/>}
    </div>
  );
};

export default SinglePlayerGame;
