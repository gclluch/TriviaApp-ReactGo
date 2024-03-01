// MainMenu.tsx
import React from 'react';
import { Link } from 'react-router-dom';

const MainMenu: React.FC = () => {
  const menuItems = [
    { path: "/singleplayer", label: "Single Player" },
    { path: "/multiplayer", label: "Multiplayer" },
  ];

  return (
    <div className="menu">
      <h1 className="welcome-heading">Welcome to CapTrivia</h1>
      {menuItems.map((item, index) => (
        <Link key={index} to={item.path}>
          <button className="menu-button">{item.label}</button>
        </Link>
      ))}
    </div>
  );
};

export default MainMenu;
