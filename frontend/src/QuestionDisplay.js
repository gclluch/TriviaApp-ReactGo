// QuestionDisplay.js
import React from 'react';
import PropTypes from 'prop-types';

const QuestionDisplay = ({ question, options, onAnswer }) => (
  <div>
    <h3>{question}</h3>
    {options ? options.map((option, index) => (
      <button key={index} onClick={onAnswer(index)} className="option-button">
        {option}
      </button>
    )) : <p>Loading options...</p>}
  </div>
);

QuestionDisplay.propTypes = {
  question: PropTypes.string.isRequired,
  options: PropTypes.arrayOf(PropTypes.string).isRequired,
  onAnswer: PropTypes.func.isRequired,
};


export default QuestionDisplay;
