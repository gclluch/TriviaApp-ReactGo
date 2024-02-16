// MainMenu.js
import React from 'react';
import { Link } from 'react-router-dom';

function MainMenu() {
  return (
    <div className="menu">
    <h1 className="welcome-heading">Welcome to CapTrivia</h1>
      <Link to="/singleplayer"><button className="menu-button">Single Player</button></Link>
      <Link to="/multiplayer"><button className="menu-button">Multiplayer</button></Link>
    </div>
  );
}

export default MainMenu;
