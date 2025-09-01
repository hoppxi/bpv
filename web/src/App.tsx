import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { PlayerProvider } from "@/hooks/usePlayer";
import { LibraryProvider } from "@/hooks/useLIbrary";
import Layout from "@/components/Layout/Layout";
import Library from "@/components/Library/Library";
import Player from "@/components/Player/Player";
import "@/styles/App.css";

function App() {
  return (
    <PlayerProvider>
      <LibraryProvider>
        <Router>
          <div className="app">
            <Layout>
              <Routes>
                <Route path="/" element={<Library />} />
                <Route path="/library" element={<Library />} />
                <Route path="/album/:id" element={<div>Album View</div>} />
                <Route path="/artist/:id" element={<div>Artist View</div>} />
              </Routes>
            </Layout>
            <Player />
          </div>
        </Router>
      </LibraryProvider>
    </PlayerProvider>
  );
}

export default App;
