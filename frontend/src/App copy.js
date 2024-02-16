import React, { useState, useEffect } from "react";
import "./App.css";
import StartGameButton from './StartGameButton';
import QuestionDisplay from './QuestionDisplay';
import ScoreDisplay from './ScoreDisplay';
import LoadingIndicator from './LoadingIndicator';
import ErrorMessage from './ErrorMessage';

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import JoinGameComponent from './JoinGameComponent'; // Make sure this component is correctly implemented
import { useWebSocket } from './WebSocketContext'; // Import the useWebSocket hook


// Use REACT_APP_BACKEND_URL or http://localhost:8080 as the API_BASE
const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

function App() {
  const [gameSession, setGameSession] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [score, setScore] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [shareableLink, setShareableLink] = useState("");
  const [theme, setTheme] = useState('dark'); // Default theme

  const toggleTheme = () => {
    setTheme((currentTheme) => (currentTheme === 'light' ? 'dark' : 'light'));
  };

  useEffect(() => {
    document.body.setAttribute('data-theme', theme);
  }, [theme]);

  const webSocket = useWebSocket();

  useEffect(() => {
    if (!webSocket) return;

    webSocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      switch (message.type) {
        case 'WAITING_ROOM_UPDATE':
          // Handle waiting room updates here
          break;
        case 'GAME_START':
          // Handle game start here
          break;
        case 'QUESTION':
          // Display the current question to players here
          break;
        case 'SCORE_UPDATE':
          // Update scores based on answers here
          break;
        case 'GAME_END':
          // Display the winner and handle game end here
          break;
        default:
          console.log('Received message from server: ', message);
      }
    };
  }, [webSocket]);

  const startGame = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = await res.json();
      setGameSession(data.sessionId);
      // setShareableLink(data.shareableLink); // Store the shareable link from the response
      console.log(data.shareableLink);
      fetchQuestions();
    } catch (err) {
      setError("Failed to start game.");
    }
    setLoading(false);
  };

  const fetchQuestions = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/questions`);
      const data = await res.json();
      setQuestions(data);
    } catch (err) {
      setError("Failed to fetch questions.");
    }
    setLoading(false);
  };

  // modified to return function that takes index as argument.
  // this way we can call submitAnswer from QuestionDisplay.js
  // without creating new function every render
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

  // consider only loading in certain cases to avoid loading icon
  if (loading) return <LoadingIndicator />;


  if (error) return <ErrorMessage message={error} />;

  return (
    <div className="App">
      <button onClick={toggleTheme} className="toggle-theme-button">
      Toggle Theme
      </button>
      <main>
        {!gameSession ? (
          <StartGameButton onStart={startGame} />
        ) : (
          <div>
            <QuestionDisplay
              question={questions[currentQuestionIndex]?.questionText}
              options={questions[currentQuestionIndex]?.options}
              onAnswer={submitAnswer}
            />
            <ScoreDisplay score={score}/>
          </div>
        )}
      </main>
    </div>
  );

  return (
    <Router>
      <div className="App">
        <button onClick={toggleTheme} className="toggle-theme-button">
          Toggle Theme
        </button>
        <Routes>
          <Route path="/" element={
            <main>
              {!gameSession ? (
                <StartGameButton onStart={startGame} />
              ) : (
                <div>
                  <QuestionDisplay
                    question={questions[currentQuestionIndex]?.questionText}
                    options={questions[currentQuestionIndex]?.options}
                    onAnswer={submitAnswer}
                  />
                  <ScoreDisplay score={score}/>
                </div>
              )}
            </main>
          } />
          <Route path="/join/:sessionId" element={<JoinGameComponent />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
