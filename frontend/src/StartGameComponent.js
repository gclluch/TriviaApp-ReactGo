import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const StartGameComponent = () => {
  const navigate = useNavigate();
  const [shareableLink, setShareableLink] = useState("");
  const [numQuestions, setNumQuestions] = useState(10); // State to hold the number of questions

  const startNewGame = async () => {
    try {
      const response = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ numQuestions }), // Send the number of questions to the backend
      });
      const data = await response.json();
      setShareableLink(`${window.location.origin}/join/${data.sessionId}`);
      navigate(`/join/${data.sessionId}`); // Navigate to the shareable link
    } catch (error) {
      console.error("Failed to start a new game:", error);
    }
  };

  const handleNumQuestionsChange = (e) => {
    const value = Math.max(1, Math.min(50, Number(e.target.value))); // Ensure value is between 1 and 20
    setNumQuestions(value);
  };

  return (
    <div>
      <h1>Multiplayer</h1>
      <div>
        <label htmlFor="numQuestions">Number of Questions: </label>
        <input
          id="numQuestions"
          type="number"
          value={numQuestions}
          onChange={handleNumQuestionsChange}
          min="1"
          max="20"
        />
      </div>
      <button onClick={startNewGame}>Create Multiplayer Session</button>
      {shareableLink && (
        <>
          <p>Shareable link generated:</p>
          <input type="text" value={shareableLink} readOnly onFocus={(e) => e.target.select()} />
        </>
      )}
    </div>
  );
};

export default StartGameComponent;
