import { useState } from "react";
import { Search, Menu, User, Settings } from "lucide-react";
import "@/styles/Layout/Header.css";

export default function Header() {
  const [searchQuery, setSearchQuery] = useState("");

  return (
    <header className="header">
      <div className="header-left">
        <button className="menu-button">
          <Menu size={24} />
        </button>
        <h1 className="header-title">BPV Music Player</h1>
      </div>

      <div className="header-center">
        <div className="search-container">
          <Search size={20} className="search-icon" />
          <input
            type="text"
            placeholder="Search songs, artists, albums..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
          />
        </div>
      </div>

      <div className="header-right">
        <button className="header-button">
          <Settings size={20} />
        </button>
        <button className="header-button">
          <User size={20} />
        </button>
      </div>
    </header>
  );
}
