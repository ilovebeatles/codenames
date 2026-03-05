import { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useGameStore } from '../store/gameStore';
import { useWebSocket } from '../ws/useWebSocket';
import { joinRoom } from '../api/http';
import Board from '../components/Board';
import ClueInput from '../components/ClueInput';
import ClueDisplay from '../components/ClueDisplay';
import GameOverBanner from '../components/GameOverBanner';
import JoinTeamModal from '../components/JoinTeamModal';
import type { Team } from '../types';

export default function GamePage() {
  const { roomID } = useParams<{ roomID: string }>();
  const navigate = useNavigate();
  const roomState = useGameStore((s) => s.roomState);
  const playerName = useGameStore((s) => s.playerName);
  const error = useGameStore((s) => s.error);
  const clearError = useGameStore((s) => s.clearError);
  const { send } = useWebSocket(roomID);

  // Ensure player is joined
  useEffect(() => {
    if (!playerName) {
      navigate(`/room/${roomID}`);
      return;
    }
    if (roomID) {
      joinRoom(roomID, playerName).catch(() => {});
    }
  }, [roomID, playerName, navigate]);

  // Navigate back to lobby if game is reset
  useEffect(() => {
    if (roomState && !roomState.game) {
      navigate(`/room/${roomID}`);
    }
  }, [roomState, roomID, navigate]);

  const game = roomState?.game;
  const cards = roomState?.cards ?? [];
  const players = roomState?.players ?? [];

  const currentPlayer = players.find((p) => p.name === playerName);
  const isSpymaster = currentPlayer?.role === 'spymaster';
  const isOperative = currentPlayer?.role === 'operative';
  const isMyTeamTurn = game?.current_team === currentPlayer?.team;
  const isCluePhase = game?.phase === 'playing' && !game.current_clue;
  const isGuessPhase = game?.phase === 'playing' && !!game.current_clue;

  const canGiveClue = isSpymaster && isMyTeamTurn && isCluePhase;
  const canGuess = isOperative && isMyTeamTurn && isGuessPhase;
  const canEndGuessing = isOperative && isMyTeamTurn && isGuessPhase;

  const teamColor = game?.current_team === 'red' ? '#d32f2f' : '#1976d2';
  const teamName = game?.current_team === 'red' ? 'Красные' : 'Синие';

  const redRemaining = roomState?.red_cards_left
  const blueRemaining = roomState?.blue_cards_left

  const needsTeam = currentPlayer && !currentPlayer.team;

  const handleJoinTeam = (team: Team) => {
    send({ type: 'join_team', team, role: 'operative' });
  };

  const handleGiveClue = (clue: string, number: number) => {
    send({ type: 'give_clue', clue, number });
  };

  const handleGuess = (cardID: string) => {
    send({ type: 'guess_card', card_id: cardID });
  };

  const handleEndGuessing = () => {
    send({ type: 'end_guessing' });
  };

  const handleNewGame = () => {
    send({ type: 'new_game' });
  };

  if (!game) return <div className="loading">Загрузка...</div>;

  return (
    <div className="game-page">
      {needsTeam && game.phase === 'playing' && (
        <JoinTeamModal players={players} onJoin={handleJoinTeam} />
      )}

      {game.phase === 'finished' && (
        <GameOverBanner winner={game.winner} onNewGame={handleNewGame} />
      )}

      {error && (
        <div className="error-banner" onClick={clearError}>
          {error}
        </div>
      )}

      <div className="game-header">
        <div className="score">
          <span className="score-red">{redRemaining}</span>
          {' - '}
          <span className="score-blue">{blueRemaining}</span>
        </div>
        {game.phase === 'playing' && (
          <div className="turn-indicator" style={{ color: teamColor }}>
            Ход: {teamName}
            {isCluePhase && ' (подсказка)'}
            {isGuessPhase && ' (угадывание)'}
          </div>
        )}
        {isSpymaster && (
          <div className="spy-badge">Вы — Спаймастер</div>
        )}
      </div>

      {isGuessPhase && game.current_clue && (
        <ClueDisplay
          clue={game.current_clue}
          number={game.current_number}
          guessesLeft={game.guesses_left}
        />
      )}

      {canGiveClue && <ClueInput onGiveClue={handleGiveClue} />}

      <Board cards={cards} onGuess={handleGuess} canGuess={canGuess} />

      {canEndGuessing && (
        <button className="end-guessing-btn" onClick={handleEndGuessing}>
          Закончить угадывание
        </button>
      )}

      <div className="game-players">
        <div className="game-team red-team">
          <h4>Красные</h4>
          {players.filter((p) => p.team === 'red').map((p) => (
            <span key={p.id} className={`player-tag ${!p.is_online ? 'offline' : ''}`}>
              {p.name} {p.role === 'spymaster' ? '🕵️' : '🔍'}
            </span>
          ))}
        </div>
        <div className="game-team blue-team">
          <h4>Синие</h4>
          {players.filter((p) => p.team === 'blue').map((p) => (
            <span key={p.id} className={`player-tag ${!p.is_online ? 'offline' : ''}`}>
              {p.name} {p.role === 'spymaster' ? '🕵️' : '🔍'}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
