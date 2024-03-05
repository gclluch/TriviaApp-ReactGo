import React, { ChangeEvent, useState } from 'react'; // 1
import { useNavigate } from 'react-router-dom';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080";

const StartGameComponent: React.FC = () => { // 2
  const navigate = useNavigate();
  const [shareableLink, setShareableLink] = useState<string>(""); // 3
  const [numQuestions, setNumQuestions] = useState<number>(10); // 4

  const startNewGame = async () => {
    try {
      const response = await fetch(`${API_BASE}/game/start`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ numQuestions }),
      });
      const data: { sessionId: string } = await response.json(); // 5
      setShareableLink(`${window.location.origin}/join/${data.sessionId}`);
      navigate(`/join/${data.sessionId}`);
    } catch (error) {
      console.error("Failed to start a new game:", error);
    }
  };

  const handleNumQuestionsChange = (e: ChangeEvent<HTMLInputElement>) => { // 6
    const value = Math.max(1, Math.min(50, Number(e.target.value)));
    setNumQuestions(value);
  };

  return (
    <div>
      <h1>Multiplayer</h1>
      <div>
        <label htmlFor="numQuestions">Number of Questions: </label>
        <input
          id="numQuestions"
          type="number"
          value={numQuestions}
          onChange={handleNumQuestionsChange}
          min="1"
          max="50" // 7
        />
      </div>
      <button onClick={startNewGame}>Create Multiplayer Session</button>
      {shareableLink && (
        <>
          <p>Shareable link generated:</p>
          <input type="text" value={shareableLink} readOnly onFocus={(e: ChangeEvent<HTMLInputElement>) => e.target.select()} /> // 8
        </>
      )}
    </div>
  );
};

export default StartGameComponent;
