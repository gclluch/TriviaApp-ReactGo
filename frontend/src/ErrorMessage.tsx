// ErrorMessage.tsx
import React from 'react';

// Define an interface for the component's props
interface ErrorMessageProps {
  message: string;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message }) => (
  <div aria-live="assertive" className="error">Error: {message}</div>
);

export default ErrorMessage;
