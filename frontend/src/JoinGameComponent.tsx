import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useWebSocket } from './WebSocketContext'; // Assume this is already correctly implemented in TypeScript

interface RouteParams {
  [key: string]: string | undefined;
}

interface JoinGameData {
  playerName: string;
  playerId: string;
}

const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080';

const JoinGameComponent: React.FC = () => {
  const { sessionId } = useParams<RouteParams>();
  const navigate = useNavigate();
  const { webSocket, isConnected } = useWebSocket();

  const [shareableLink, setShareableLink] = useState<string>('');
  const [hasJoined, setHasJoined] = useState<boolean>(false);
  const [playerName, setPlayerName] = useState<string>('');
  const [playerId, setPlayerId] = useState<string>('');
  const [playerCount, setPlayerCount] = useState<number>(0);
  const [countdown, setCountdown] = useState<number | null>(null);

  const joinGame = async () => {
    console.log(`Attempting to join game session: ${sessionId}`);
    try {
      const response = await fetch(`${API_BASE}/game/join/${sessionId}`, {
        method: 'POST',
      });

      if (!response.ok) throw new Error('Network response was not ok.');
      const data: JoinGameData = await response.json();
      console.log('Successfully joined the game:', data);
      setHasJoined(true);
      setPlayerName(data.playerName);
      setPlayerId(data.playerId);
    } catch (error) {
      console.error('Failed to join the game:', error);
    }
  };

  useEffect(() => {
    setShareableLink(`${window.location.origin}/join/${sessionId}`);
  }, [sessionId]);

  useEffect(() => {
    if (webSocket && isConnected) {
      const handleMessage = (event: MessageEvent) => {
        const data = JSON.parse(event.data);
        switch (data.type) {
          case 'playerCount':
            setPlayerCount(data.count);
            break;
          case 'countdown':
            setCountdown(data.time);
            break;
          default:
            console.log('Unhandled message type:', data.type);
        }
      };

      webSocket.addEventListener('message', handleMessage);

      if (webSocket.readyState === WebSocket.OPEN) {
        webSocket.send(JSON.stringify({ action: 'joinSession', sessionId }));
      } else {
        webSocket.onopen = () => {
          console.log('WebSocket Connected');
          webSocket.send(JSON.stringify({ action: 'joinSession', sessionId }));
        };
      }

      return () => webSocket.removeEventListener('message', handleMessage);
    }
  }, [webSocket, isConnected, sessionId]);

  useEffect(() => {
    if (countdown === 0 && hasJoined) {
      navigate(`/game/${sessionId}`, {
        state: { playerName, playerId, gameStarted: true },
      });
    } else if (countdown === 0) {
      navigate(`/game/${sessionId}`, { state: { gameStarted: false } });
    }
  }, [countdown, navigate, hasJoined, playerName, sessionId, playerId]);

  return (
    <div>
      {hasJoined ? (
        <h3>Joined Game as: {playerName}</h3>
      ) : (
        <h3>Joining Game Session...</h3>
      )}
      <p>Share this link to invite more players:</p>
      <input
        type="text"
        value={shareableLink}
        readOnly
        onFocus={e => e.target.select()}
      />
      {!hasJoined && (
        <div style={{ marginTop: '10px' }}>
          <button onClick={() => joinGame()}>Join Game</button>
        </div>
      )}
      <div className="footer">
        {countdown !== null && <div>Game Starting in: {countdown}</div>}
        <h4>Player Count: {playerCount}</h4>
      </div>
    </div>
  );
};

export default JoinGameComponent;
