// ErrorMessage.js
import React from 'react';
import PropTypes from 'prop-types';

const ErrorMessage = ({ message }) => (
  <div aria-live="assertive" className="error">Error: {message}</div>
);

ErrorMessage.propTypes = {
  message: PropTypes.string.isRequired,
};

export default ErrorMessage;
