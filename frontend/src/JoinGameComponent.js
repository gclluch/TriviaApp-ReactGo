import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useWebSocket } from './WebSocketContext'; // Ensure this is correctly implemented

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const JoinGameComponent = () => {
  const { sessionId } = useParams();
  const [shareableLink, setShareableLink] = useState("");
  const [hasJoined, setHasJoined] = useState(false);
  const [playerName, setPlayerName] = useState("");
  const [playerCount, setPlayerCount] = useState(0); // State for tracking live player count
  const { webSocket, isConnected } = useWebSocket(); // Destructuring to get webSocket and isConnected
  const [countdown, setCountdown] = useState(null);
  const navigate = useNavigate();


  const joinGame = async (sessionId) => {
    console.log(`Attempting to join game session: ${sessionId}`);
    try {
      const response = await fetch(`${API_BASE}/game/join`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sessionId: sessionId }),
      });
      if (!response.ok) throw new Error('Network response was not ok.');
      const data = await response.json();
      console.log("Successfully joined the game:", data);
      setHasJoined(true);
      setPlayerName(data.playerName); // Update player name upon joining
    } catch (error) {
      console.error("Failed to join the game:", error);
    }
  };

  useEffect(() => {
    setShareableLink(`${window.location.origin}/join/${sessionId}`);
  }, [sessionId]);

  useEffect(() => {
    if (webSocket && isConnected) {
      const handleMessage = (event) => {
        const data = JSON.parse(event.data);
        // console.log('Received message:', data);
        switch (data.type) {
          case 'playerCount':
            setPlayerCount(data.count);
            break;
          case 'countdown':
            setCountdown(data.time);
            console.log("Countdown:", data.time);
            break;
          default:
            console.log("Unhandled message type:", data.type);
        }
      };

      webSocket.addEventListener('message', handleMessage);

      // Send joinSession message once WebSocket is connected
      if (webSocket.readyState === WebSocket.OPEN) {
        webSocket.send(JSON.stringify({ action: 'joinSession', sessionId }));
      } else {
        webSocket.onopen = () => {
          console.log('WebSocket Connected');
          webSocket.send(JSON.stringify({ action: 'joinSession', sessionId }));
        };
      }

      return () => {
        webSocket.removeEventListener('message', handleMessage);
      };
    }
  }, [webSocket, isConnected, sessionId]);

  useEffect(() => {
    console.log("Player count:", playerCount);
    if (countdown === 0 && hasJoined) {

      navigate(`/game/${sessionId}`, { state: { playerName: playerName } }); // Passing player name in state
    }
  }, [countdown, navigate, hasJoined, playerName, sessionId]);

  return (
    <div>
      {hasJoined ? <h3>Joined Game as: {playerName}</h3> : <h3>Joining Game Session...</h3>}
      <p>Share this link to invite more players:</p>
      <input type="text" value={shareableLink} readOnly onFocus={(e) => e.target.select()} />
      {!hasJoined && (
        <div style={{marginTop: "10px"}}>
          <button onClick={() => joinGame(sessionId)}>Join Game</button>
        </div>
      )}
      {/* Display live player count */}
      <div className="footer"> {/* Use the footer class for player count */}
        <h4>Live Player Count: {playerCount}</h4>
      </div>
      {countdown !== null && <div>Countdown: {countdown}</div>}
          </div>
  );
};

export default JoinGameComponent;
