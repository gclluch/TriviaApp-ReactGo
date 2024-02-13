import React, { createContext, useContext, useEffect, useState } from 'react';

const WebSocketContext = createContext(null);

export const useWebSocket = () => useContext(WebSocketContext);

export const WebSocketProvider = ({ children }) => {
    const [webSocket, setWebSocket] = useState(null);

    useEffect(() => {
        let ws;
        const connect = () => {
            ws = new WebSocket('ws://localhost:8080/ws');

            ws.onopen = () => console.log('WebSocket Connected');
            ws.onclose = () => {
                console.log('WebSocket Disconnected');
                // setTimeout(connect, 3000); // Attempt to reconnect every 3 seconds
            };
            ws.onerror = (error) => console.log('WebSocket Error:', error);

            setWebSocket(ws);
        };

        connect();

        return () => {
            if (ws) {
                ws.close();
            }
        };
    }, []);

    return (
        <WebSocketContext.Provider value={webSocket}>
            {children}
        </WebSocketContext.Provider>
    );
};
