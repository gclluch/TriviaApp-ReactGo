import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import "./App.css";
import MainMenu from './MainMenu';
import SinglePlayerGame from './SinglePlayerGame';
import StartGameComponent from './StartGameComponent';
import JoinGameComponent from './JoinGameComponent';
import MultiplayerGame from "./MultiplayerGame";
import FinalScores from "./FinalScores";
import Leaderboard from "./Leaderboard";

const App: React.FC = () => {
  const [theme, setTheme] = useState<string>('dark');

  useEffect(() => {
    document.body.setAttribute('data-theme', theme);
  }, [theme]);

  const toggleTheme = (): void => {
    setTheme(current => current === 'light' ? 'dark' : 'light');
  };

  return (
    <Router>
      <div className="App">
        <header className="app-header">
          <button onClick={toggleTheme} className="header-button">Toggle Theme</button>
          <Link to="/"><button className="main-menu-button">Main Menu</button></Link>
          <Link to="/leaderboard"><button className="header-button">Leaderboard</button></Link>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<MainMenu />} />
            <Route path="/singleplayer" element={<SinglePlayerGame />} />
            <Route path="/multiplayer" element={<StartGameComponent />} />
            <Route path="/join/:sessionId" element={<JoinGameComponent />} />
            <Route path="/game/:sessionId" element={<MultiplayerGame />} />
            <Route path="/final-scores/:sessionId" element={<FinalScores />} />
            <Route path="/leaderboard" element={<Leaderboard />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
