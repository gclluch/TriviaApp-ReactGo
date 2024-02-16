import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useWebSocket } from './WebSocketContext';

const JoinGameComponent = () => {
    const { sessionId } = useParams();
    const navigate = useNavigate();
    const webSocket = useWebSocket();

    useEffect(() => {
      if (!webSocket) return;

      const joinGame = () => {
          if (webSocket.readyState === WebSocket.OPEN) {
              webSocket.send(JSON.stringify({ action: 'join', sessionId }));
          } else {
              webSocket.addEventListener('open', () => {
                  webSocket.send(JSON.stringify({ action: 'join', sessionId }));
              }, { once: true });
          }
      };

      joinGame();
  }, [sessionId, webSocket]);

    return <div>Attempting to join game session...</div>;
};

export default JoinGameComponent;
