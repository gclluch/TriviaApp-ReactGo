// WebSocketContext.tsx
import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';

// Define an interface for the context value
interface WebSocketContextValue {
  webSocket: WebSocket | null;
  isConnected: boolean;
}

// Define an interface for the WebSocketProvider props
interface WebSocketProviderProps {
  children: ReactNode;
  url?: string;
}

// Create context with an initial value
const WebSocketContext = createContext<WebSocketContextValue>({
  webSocket: null,
  isConnected: false,
});

// Custom hook to use the WebSocket context
export const useWebSocket = (): WebSocketContextValue => useContext(WebSocketContext);

// WebSocketProvider component
export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children, url = 'ws://localhost:8080/ws' }) => {
  const [webSocket, setWebSocket] = useState<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState<boolean>(false);

  useEffect(() => {
    const ws = new WebSocket(url);

    const onOpen = () => {
      console.log('WebSocket Connected');
      setIsConnected(true);
    };

    const onMessage = (event: MessageEvent) => {
      console.log('Received message:', event.data);
    };

    const onClose = () => {
      console.log('WebSocket Disconnected');
      setIsConnected(false);
    };

    const onError = (error: Event) => {
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
