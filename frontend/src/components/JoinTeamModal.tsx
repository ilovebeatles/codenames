import type { Team, Player } from '../types';

interface Props {
  players: Player[];
  onJoin: (team: Team) => void;
}

export default function JoinTeamModal({ players, onJoin }: Props) {
  const redPlayers = players.filter((p) => p.team === 'red');
  const bluePlayers = players.filter((p) => p.team === 'blue');

  return (
    <div className="modal-overlay">
      <div className="modal join-team-modal">
        <h2>Выберите команду</h2>
        <div className="join-team-options">
          <button
            className="join-team-btn join-team-red"
            onClick={() => onJoin('red')}
          >
            <span className="join-team-name">Красные</span>
            <span className="join-team-count">{redPlayers.length} игроков</span>
          </button>
          <button
            className="join-team-btn join-team-blue"
            onClick={() => onJoin('blue')}
          >
            <span className="join-team-name">Синие</span>
            <span className="join-team-count">{bluePlayers.length} игроков</span>
          </button>
        </div>
      </div>
    </div>
  );
}
