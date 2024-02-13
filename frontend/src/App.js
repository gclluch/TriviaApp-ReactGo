import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import "./App.css";
import MainMenu from './MainMenu';
import SinglePlayerGame from './SinglePlayerGame';
import MultiplayerGame from './MultiplayerGame';

function App() {
  const [theme, setTheme] = useState('dark');

  useEffect(() => {
    document.body.setAttribute('data-theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme(current => current === 'light' ? 'dark' : 'light');
  };

  return (
    <Router>
      <div className="App">
        <header className="app-header">
          <button onClick={toggleTheme} className="header-button">Toggle Theme</button>
          <Link to="/"><button className="main-menu-button">Main Menu</button></Link>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<MainMenu />} />
            <Route path="/singleplayer" element={<SinglePlayerGame />} />
            <Route path="/multiplayer" element={<MultiplayerGame />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
