// MultiplayerGame.js
import React, { useState, useEffect } from "react";
import { useNavigate } from 'react-router-dom';
import { useWebSocket } from './WebSocketContext'; // Ensure this hook is implemented correctly

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const MultiplayerGame = () => {
    const [gameSession, setGameSession] = useState(null);
    const [shareableLink, setShareableLink] = useState("");
    const navigate = useNavigate();
    const webSocket = useWebSocket(); // Use the WebSocket connection provided by the context

    useEffect(() => {
      if (webSocket) {
          webSocket.onmessage = (event) => {
              const data = JSON.parse(event.data);
              console.log("Received data:", data);
              // Handle the received data accordingly
          };

          webSocket.onopen = () => {
              console.log("WebSocket connection established");
              // If joining an existing session, send a join message
              if (gameSession) {
                  webSocket.send(JSON.stringify({ action: 'join', sessionId: gameSession }));
              }
          };
      }

      return () => {
          if (webSocket) {
              webSocket.onmessage = null;
              webSocket.onopen = null;
          }
      };
  }, [webSocket, gameSession]); // Re-run effect if webSocket or gameSession changes

  
    // Function to start a new game session
    const startGame = async () => {
        try {
            const response = await fetch(`${API_BASE}/game/start`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
            });
            const data = await response.json();
            setGameSession(data.sessionId);
            setShareableLink(`${window.location.origin}/join/${data.sessionId}`);
            // You can also establish WebSocket connection here if it's not globally available
        } catch (error) {
            console.error("Failed to start multiplayer game:", error);
        }
    };

    // This effect runs once on component mount to start a new game session
    useEffect(() => {
        startGame();
    }, []); // Empty dependency array ensures this effect runs only once on mount

    // JSX for the component
    return (
        <div>
            {!gameSession ? (
                <p>Creating a new multiplayer session...</p>
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
