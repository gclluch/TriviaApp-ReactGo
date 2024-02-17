import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from 'react-router-dom';
// import { useWebSocket } from './WebSocketContext'; // Ensure this hook is implemented correctly

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
  const { sessionId } = useParams(); // Extract sessionId from the route
  console.log(sessionId);
  const navigate = useNavigate();
  // const webSocket = useWebSocket(); // Use the WebSocket connection provided by the context
  const [shareableLink, setShareableLink] = useState("");

  // setShareableLink(`${window.location.origin}/join/${sessionId}`);

  useEffect(() => {
    if (!sessionId) {
      startNewGame();
    } else {
      joinGame(sessionId);
    }
  }, [sessionId]);


  const startNewGame = async () => {
    try {
      const response = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });
      if (!response.ok) throw new Error('Network response was not ok.');
      const data = await response.json();
      setShareableLink(data.shareableLink); // Set the shareable link for display
      navigate(`/join/${data.sessionId}`); // Navigate to the shareable link
    } catch (error) {
      console.error("Failed to start a new game:", error);
    }
  };

  const joinGame = async (sessionId) => {
    // Here, you might call an endpoint like `/game/join` to add the player to the session.
    // This example doesn't directly call such an endpoint but it's where you would do so.
    console.log(`Joining game session: ${sessionId}`);
    // Optionally, set the shareable link based on the sessionId for consistency
    setShareableLink(`${window.location.origin}/join/${sessionId}`);
  };

  return (
    <div>
      <h2>Multiplayer Session</h2>
      {shareableLink && (
        <>
          <p>Share this link to invite more players:</p>
          <input type="text" value={shareableLink} readOnly onFocus={(e) => e.target.select()} />
        </>
      )}
      <button onClick={() => navigate('/')}>Back to Main Menu</button>
    </div>
  );
};

export default MultiplayerGame;