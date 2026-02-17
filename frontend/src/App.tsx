import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import LobbyPage from './pages/LobbyPage';
import GamePage from './pages/GamePage';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/room/:roomID" element={<LobbyPage />} />
        <Route path="/room/:roomID/game" element={<GamePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
