import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useWebSocket } from './WebSocketContext'; // Import your hook

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const JoinGameComponent = () => {
  const { sessionId } = useParams();
  const [shareableLink, setShareableLink] = useState("");
  const [hasJoined, setHasJoined] = useState(false); // New state to track if the user has joined
  const [playerName, setPlayerName] = useState(""); // Store the player name
  const [playerCount, setPlayerCount] = useState(0); // State to keep track of the live player count
  const { webSocket, isConnected } = useWebSocket(); // Destructure to get both webSocket and isConnected


  const joinGame = async (sessionId) => {
    console.log(`Attempting to join game session: ${sessionId}`);
    try {
      const response = await fetch(`${API_BASE}/game/join`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          sessionId: sessionId,
        }),
      });
      if (!response.ok) throw new Error('Network response was not ok.');
      const data = await response.json();
      console.log("Successfully joined the game:", data);
      setHasJoined(true); // Update state to indicate the user has joined
      setPlayerName(data.playerName); // Set the player name upon joining

    } catch (error) {
      console.error("Failed to join the game:", error);
    }
  };

  useState(() => {
    setShareableLink(`${window.location.origin}/join/${sessionId}`);
  }, [sessionId]);


  // WEBSOCKET

  useEffect(() => {
    if (webSocket && isConnected) {
      // Define the message handler
      const handleMessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'playerCount') {
          console.log('Received player count:', data.count);
          setPlayerCount(data.count); // Update state with the new player count
        }
        // Handle other message types as needed
      };

      webSocket.addEventListener('message', handleMessage);
      console.log('WebSocket Connected - Sending joinSession message');
      webSocket.send(JSON.stringify({ action: 'joinSession', sessionId }));


      // Wait for the connection to open before sending a message
      webSocket.onopen = () => {
        console.log('WebSocket Connected');
        webSocket.send(JSON.stringify({ action: 'joinSession', sessionId }));
      };

      // Clean up the event listener and onopen handler when the component unmounts or websocket changes
      return () => {
        webSocket.removeEventListener('message', handleMessage);
        webSocket.onclose = () => console.log('WebSocket Disconnected'); // Optionally set onclose handler
      };
    }
  }, [webSocket, isConnected, sessionId]); // Re-run the effect if webSocket or sessionId changes


  return (
    <div>
      {hasJoined ? <h3>Joined Game as: {playerName}</h3> : <h3>Joining Game Session...</h3>}
      <p>Share this link to invite more players:</p>
      <input type="text" value={shareableLink} readOnly onFocus={(e) => e.target.select()} />
    {/* Moved the button below the input field */}
    {!hasJoined && (
      <div style={{marginTop: "10px"}}>
        <button onClick={() => joinGame(sessionId)}>Join Game</button>
      </div>
    )}
  </div>
  );
};

export default JoinGameComponent;
