import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createRoom } from '../api/http';
import { useGameStore } from '../store/gameStore';
import PlayerNameModal from '../components/PlayerNameModal';

export default function LandingPage() {
  const navigate = useNavigate();
  const playerName = useGameStore((s) => s.playerName);
  const setPlayerName = useGameStore((s) => s.setPlayerName);
  const [showNameModal, setShowNameModal] = useState(!playerName);

  const handleCreate = async () => {
    if (!playerName) {
      setShowNameModal(true);
      return;
    }
    const room = await createRoom();
    navigate(`/room/${room.id}`);
  };

  const handleNameSubmit = (name: string) => {
    setPlayerName(name);
    setShowNameModal(false);
  };

  return (
    <div className="landing">
      {showNameModal && <PlayerNameModal onSubmit={handleNameSubmit} />}
      <h1>Codenames Online</h1>
      <p>Онлайн-версия настольной игры Codenames</p>
      {playerName && <p>Привет, {playerName}!</p>}
      <button className="create-btn" onClick={handleCreate}>
        Создать комнату
      </button>
      {playerName && (
        <button className="change-name-btn" onClick={() => setShowNameModal(true)}>
          Сменить имя
        </button>
      )}
    </div>
  );
}
