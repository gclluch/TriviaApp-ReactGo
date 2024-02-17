import React, { createContext, useContext, useState, useEffect } from 'react';

// Create a WebSocket context
const WebSocketContext = createContext(null);

// Custom hook for components to access the WebSocket context
export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (context === undefined) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
};

// Provider component that wraps your app and provides the WebSocket context
export const WebSocketProvider = ({ children }) => {
    const [webSocket, setWebSocket] = useState(null);
    const [isConnected, setIsConnected] = useState(false); // State to track connection status

    useEffect(() => {
        // Initialize WebSocket connection
        const ws = new WebSocket('ws://localhost:8080/ws');

        // Set up WebSocket event listeners
        ws.onopen = () => {
            console.log('WebSocket Connected');
            setIsConnected(true); // Update connection status
        };

        ws.onmessage = (event) => {
            console.log('Received message:', event.data);
            // Here you can handle incoming messages
        };

        ws.onclose = () => {
            console.log('WebSocket Disconnected');
            setIsConnected(false); // Update connection status
        };

        // Update the WebSocket state
        setWebSocket(ws);

        // Clean up on component unmount
        return () => {
            ws.close();
        };
    }, []);

    // Provide the WebSocket and its connection status to the context consumers
    return (
        <WebSocketContext.Provider value={{ webSocket, isConnected }}>
            {children}
        </WebSocketContext.Provider>
    );
};
