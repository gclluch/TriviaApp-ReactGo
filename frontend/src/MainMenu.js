// MainMenu.js
import React from 'react';
import { Link } from 'react-router-dom';

function MainMenu() {
  return (
    <div className="menu">
    <h1 className="welcome-heading">Welcome to CapTrivia</h1>
      <Link to="/singleplayer"><button>Single Player</button></Link>
      <Link to="/multiplayer"><button>Start Multiplayer Game</button></Link>
    </div>
  );
}

export default MainMenu;
