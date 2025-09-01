import { NavLink } from "react-router-dom";
import { Home, Library, Search, Radio, Heart, Clock } from "lucide-react";
import "@/styles/Layout/Sidebar.css";

export default function Sidebar() {
  const navItems = [
    { path: "/", label: "Home", icon: Home },
    { path: "/library", label: "Library", icon: Library },
    { path: "/search", label: "Search", icon: Search },
    { path: "/radio", label: "Radio", icon: Radio },
    { path: "/favorites", label: "Favorites", icon: Heart },
    { path: "/recent", label: "Recently Played", icon: Clock },
  ];

  return (
    <aside className="sidebar">
      <nav className="sidebar-nav">
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `nav-item ${isActive ? "nav-item-active" : ""}`
              }
            >
              <Icon size={20} />
              <span>{item.label}</span>
            </NavLink>
          );
        })}
      </nav>

      <div className="sidebar-section">
        <h3 className="sidebar-title">Your Library</h3>
        <div className="library-stats">
          <div className="stat">
            <span className="stat-number">0</span>
            <span className="stat-label">Songs</span>
          </div>
          <div className="stat">
            <span className="stat-number">0</span>
            <span className="stat-label">Artists</span>
          </div>
          <div className="stat">
            <span className="stat-number">0</span>
            <span className="stat-label">Albums</span>
          </div>
        </div>
      </div>
    </aside>
  );
}
