// MultiplayerGame.js
import React, { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom';
import StartGameButton from './StartGameButton';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
  const [gameSession, setGameSession] = useState(null);
  const [shareableLink, setShareableLink] = useState("");
  const navigate = useNavigate();

  const startGame = async () => {
    try {
      const response = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      const data = await response.json();
      setGameSession(data.sessionId);
      setShareableLink(`${window.location.origin}/join/${data.sessionId}`);
    } catch (error) {
      console.error("Failed to start multiplayer game:", error);
    }
  };

  useEffect(() => {
    // Automatically start a multiplayer session when component mounts
    startGame();
  }, []);

  return (
    <div>
      {!gameSession ? (
        <StartGameButton onStart={startGame} />
      ) : (
        <>
          <h2>Multiplayer Session Created</h2>
          <p>Share this link to invite players:</p>
          <input type="text" value={shareableLink} readOnly />
          <button onClick={() => navigate('/')}>Back to Main Menu</button>
        </>
      )}
    </div>
  );
};

export default MultiplayerGame;
