import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const StartGameComponent = () => {
  const navigate = useNavigate();
  const [shareableLink, setShareableLink] = useState("");

  const startNewGame = async () => {
    try {
      const response = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      const data = await response.json();
      setShareableLink(`${window.location.origin}/join/${data.sessionId}`);
      navigate(`/join/${data.sessionId}`); // Navigate to the shareable link
    } catch (error) {
      console.error("Failed to start a new game:", error);
    }
  };

  return (
    <div>
      <h1>Multiplayer</h1>
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
