import React, { createContext, useContext, useEffect, useState } from 'react';

const WebSocketContext = createContext({
  webSocket: null,
  isConnected: false,
});

export const useWebSocket = () => useContext(WebSocketContext);

export const WebSocketProvider = ({ children, url = 'ws://localhost:8080/ws' }) => {
  const [webSocket, setWebSocket] = useState(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(url);

    const onOpen = () => {
      console.log('WebSocket Connected');
      setIsConnected(true);
    };

    const onMessage = (event) => {
      console.log('Received message:', event.data);
    };

    const onClose = () => {
      console.log('WebSocket Disconnected');
      setIsConnected(false);
    };

    const onError = (error) => {
      console.error('WebSocket Error:', error);
      setIsConnected(false);
    };

    ws.addEventListener('open', onOpen);
    ws.addEventListener('message', onMessage);
    ws.addEventListener('close', onClose);
    ws.addEventListener('error', onError);

    setWebSocket(ws);

    // Cleanup function to be called when the component unmounts
    return () => {
      ws.removeEventListener('open', onOpen);
      ws.removeEventListener('message', onMessage);
      ws.removeEventListener('close', onClose);
      ws.removeEventListener('error', onError);
      ws.close(); // Ensure WebSocket is closed when the context provider unmounts
    };
  }, [url]); // Reconnect WebSocket if URL changes

  return (
    <WebSocketContext.Provider value={{ webSocket, isConnected }}>
      {children}
    </WebSocketContext.Provider>
  );
};
