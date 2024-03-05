import React, { useState, useCallback } from 'react';
import QuestionDisplay from './QuestionDisplay';
import ScoreDisplay from './ScoreDisplay';
import LoadingIndicator from './LoadingIndicator';
import ErrorMessage from './ErrorMessage';
import StartGameButton from './StartGameButton';

interface Question {
  id: string; // Adjust according to actual structure
  questionText: string;
  options: string[];
}

interface GameStartResponse {
  sessionId: string;
}

interface AnswerResponse {
  correct: boolean;
}

const API_BASE: string =
  process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080';

const SinglePlayerGame: React.FC = () => {
  const [gameSession, setGameSession] = useState<string | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState<number>(0);
  const [score, setScore] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const startGame = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const startResponse = await fetch(`${API_BASE}/game/start`, {
        method: 'POST',
      });
      const startData: GameStartResponse = await startResponse.json();
      setGameSession(startData.sessionId);

      const questionsResponse = await fetch(
        `${API_BASE}/questions/${startData.sessionId}`
      );
      const questionsData = await questionsResponse.json();
      setQuestions(questionsData.questions);
    } catch (error) {
      setError('Failed to start or fetch questions for the game.');
    } finally {
      setLoading(false);
    }
  }, []);

  const submitAnswer = useCallback(
    async (index: number) => {
      const currentQuestion = questions[currentQuestionIndex];
      try {
        const answerResponse = await fetch(`${API_BASE}/answer`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            sessionId: gameSession,
            questionId: currentQuestion.id,
            answer: index,
          }),
        });
        const answerData: AnswerResponse = await answerResponse.json();

        console.log(answerData)
        if (answerData.correct) {
          setScore(prevScore => prevScore + 10); // Assuming each correct answer adds one point
        }
        if (currentQuestionIndex < questions.length - 1) {
          setCurrentQuestionIndex(prevIndex => prevIndex + 1);
        } else {
          await endGame();
        }
      } catch (err) {
        setError('Failed to submit answer.');
      }
    },
    [gameSession, questions, currentQuestionIndex]
  );

  const endGame = useCallback(async () => {
    try {
      await fetch(`${API_BASE}/game/end/${gameSession}`, { method: 'POST' });
      alert(`Game over! Your score: ${score}`);
      resetGame();
    } catch (err) {
      setError('Failed to end the game.');
    }
  }, [gameSession, score]);

  const resetGame = useCallback(() => {
    setGameSession(null);
    setQuestions([]);
    setCurrentQuestionIndex(0);
    setScore(0);
  }, []);

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage message={error} />;

  return (
    <div className="single-player-game">
      <h1>Single Player</h1>
      {!gameSession ? (
        <StartGameButton onStart={startGame} />
      ) : (
        <>
          {currentQuestionIndex < questions.length ? (
            <QuestionDisplay
              question={questions[currentQuestionIndex]?.questionText}
              options={questions[currentQuestionIndex]?.options}
              onAnswer={submitAnswer}
            />
          ) : (
            <p>Game Over. Your score: {score}.</p>
          )}
          <ScoreDisplay score={score} />
          <button onClick={resetGame} className="restart-button">
            Restart
          </button>
        </>
      )}
    </div>
  );
};

export default SinglePlayerGame;
