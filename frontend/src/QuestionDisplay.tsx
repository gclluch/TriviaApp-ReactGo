// QuestionDisplay.tsx
import React from 'react';

// Define an interface for the component props
interface QuestionDisplayProps {
  question: string;
  options: string[];
  onAnswer: (index: number) => void; // Assuming onAnswer takes the index of the option
}

const QuestionDisplay: React.FC<QuestionDisplayProps> = ({ question, options, onAnswer }) => (
  <div>
    <h3>{question}</h3>
    {options.length > 0 ? options.map((option, index) => (
      <button key={index} onClick={() => onAnswer(index)} className="option-button">
        {option}
      </button>
    )) : <p>Loading options...</p>}
  </div>
);

export default QuestionDisplay;
