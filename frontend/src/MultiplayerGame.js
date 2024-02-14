// MultiplayerGame.js
import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from 'react-router-dom';
import { useWebSocket } from './WebSocketContext'; // Ensure this hook is implemented correctly

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
  const { sessionId } = useParams(); // Extract sessionId from the route
  const navigate = useNavigate();
  const webSocket = useWebSocket(); // Use the WebSocket connection provided by the context
  const [shareableLink, setShareableLink] = useState("");


  return (
    <div>
      {sessionId ? (
        <>
          <h2>Joining Multiplayer Session</h2>
          {shareableLink && (
            <>
              <p>Share this link to invite more players:</p>
              <input type="text" value={shareableLink} readOnly onFocus={(e) => e.target.select()} />
            </>
          )}
        </>
      ) : (
        <p>Creating a new multiplayer session...</p>
      )}
      <button onClick={() => navigate('/')}>Back to Main Menu</button>
    </div>
  );
};

export default MultiplayerGame;
