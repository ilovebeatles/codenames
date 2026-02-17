import type { Team } from '../types';

interface Props {
  winner: Team;
  onNewGame: () => void;
}

export default function GameOverBanner({ winner, onNewGame }: Props) {
  const teamName = winner === 'red' ? 'Красные' : 'Синие';
  const color = winner === 'red' ? '#d32f2f' : '#1976d2';

  return (
    <div className="game-over-banner" style={{ borderColor: color }}>
      <h2 style={{ color }}>Победа: {teamName}!</h2>
      <button onClick={onNewGame}>Новая игра</button>
    </div>
  );
}
