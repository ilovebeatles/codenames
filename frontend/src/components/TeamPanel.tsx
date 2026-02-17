import type { Player, Team } from '../types';

interface Props {
  team: Team;
  players: Player[];
  currentPlayerID: string | null;
  onJoin: (team: Team) => void;
  onSetRole: (role: 'spymaster' | 'operative') => void;
}

export default function TeamPanel({ team, players, currentPlayerID, onJoin, onSetRole }: Props) {
  const teamPlayers = players.filter((p) => p.team === team);
  const isOnTeam = teamPlayers.some((p) => p.id === currentPlayerID);
  const color = team === 'red' ? '#d32f2f' : '#1976d2';

  return (
    <div className="team-panel" style={{ borderColor: color }}>
      <h3 style={{ color }}>
        {team === 'red' ? 'Красные' : 'Синие'}
      </h3>
      <div className="team-players">
        {teamPlayers.map((p) => (
          <div key={p.id} className={`team-player ${!p.is_online ? 'offline' : ''}`}>
            <span>{p.name}</span>
            <span className="role-badge">
              {p.role === 'spymaster' ? '🕵️' : p.role === 'operative' ? '🔍' : ''}
            </span>
          </div>
        ))}
      </div>
      {!isOnTeam ? (
        <button
          className="join-btn"
          style={{ backgroundColor: color }}
          onClick={() => onJoin(team)}
        >
          Присоединиться
        </button>
      ) : (
        <div className="role-buttons">
          <button onClick={() => onSetRole('spymaster')}>Спаймастер</button>
          <button onClick={() => onSetRole('operative')}>Оперативник</button>
        </div>
      )}
    </div>
  );
}
