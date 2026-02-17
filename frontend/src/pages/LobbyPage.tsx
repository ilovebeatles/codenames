import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { joinRoom } from '../api/http';
import { useGameStore } from '../store/gameStore';
import { useWebSocket } from '../ws/useWebSocket';
import TeamPanel from '../components/TeamPanel';
import PlayerNameModal from '../components/PlayerNameModal';
import type { Team } from '../types';

export default function LobbyPage() {
  const { roomID } = useParams<{ roomID: string }>();
  const navigate = useNavigate();
  const roomState = useGameStore((s) => s.roomState);
  const playerName = useGameStore((s) => s.playerName);
  const setPlayerName = useGameStore((s) => s.setPlayerName);
  const error = useGameStore((s) => s.error);
  const clearError = useGameStore((s) => s.clearError);
  const [joined, setJoined] = useState(false);
  const [showNameModal, setShowNameModal] = useState(false);
  const { send } = useWebSocket(joined ? roomID : undefined);

  useEffect(() => {
    if (!playerName) {
      setShowNameModal(true);
      return;
    }
    if (!roomID || joined) return;
    joinRoom(roomID, playerName).then(() => setJoined(true)).catch(() => {
      alert('Failed to join room');
    });
  }, [roomID, playerName, joined]);

  // Navigate to game when phase changes to playing
  useEffect(() => {
    if (roomState?.game?.phase === 'playing' || roomState?.game?.phase === 'finished') {
      navigate(`/room/${roomID}/game`);
    }
  }, [roomState?.game?.phase, roomID, navigate]);

  const handleNameSubmit = (name: string) => {
    setPlayerName(name);
    setShowNameModal(false);
  };

  const handleJoinTeam = (team: Team) => {
    send({ type: 'join_team', team, role: 'operative' });
  };

  const handleSetRole = (role: 'spymaster' | 'operative') => {
    send({ type: 'set_role', role });
  };

  const handleStart = () => {
    send({ type: 'start_game' });
  };

  const currentPlayer = roomState?.players.find(
    (p) => p.name === playerName
  );

  const inviteLink = `${window.location.origin}/room/${roomID}`;

  const copyLink = () => {
    navigator.clipboard.writeText(inviteLink);
  };

  return (
    <div className="lobby">
      {showNameModal && <PlayerNameModal onSubmit={handleNameSubmit} />}
      <h2>Комната: {roomID}</h2>

      <div className="invite-link">
        <input type="text" value={inviteLink} readOnly />
        <button onClick={copyLink}>Копировать</button>
      </div>

      {error && (
        <div className="error-banner" onClick={clearError}>
          {error}
        </div>
      )}

      <div className="teams-container">
        <TeamPanel
          team="red"
          players={roomState?.players ?? []}
          currentPlayerID={currentPlayer?.id ?? null}
          onJoin={handleJoinTeam}
          onSetRole={handleSetRole}
        />
        <TeamPanel
          team="blue"
          players={roomState?.players ?? []}
          currentPlayerID={currentPlayer?.id ?? null}
          onJoin={handleJoinTeam}
          onSetRole={handleSetRole}
        />
      </div>

      <div className="lobby-players">
        <h3>Без команды:</h3>
        {(roomState?.players ?? [])
          .filter((p) => !p.team)
          .map((p) => (
            <span key={p.id} className="player-chip">{p.name}</span>
          ))}
      </div>

      <button className="start-btn" onClick={handleStart}>
        Начать игру
      </button>
    </div>
  );
}
