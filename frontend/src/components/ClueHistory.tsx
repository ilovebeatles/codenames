import type { ClueHistoryEntry, Team } from '../types';

interface Props {
  team: Team;
  entries: ClueHistoryEntry[];
}

export default function ClueHistory({ team, entries }: Props) {
  const title = `История ${team === 'red' ? 'красных' : 'синих'}`;
  const emptyText = 'Подсказок пока не было';

  return (
    <section className={`clue-history clue-history-${team}`}>
      <h3>{title}</h3>
      {entries.length === 0 ? (
        <p className="clue-history-empty">{emptyText}</p>
      ) : (
        <ol className="clue-history-list">
          {entries.map((entry) => (
            <li key={entry.id} className="clue-history-item">
              <span className="clue-history-word">{entry.clue}</span>
              <span className="clue-history-number">{entry.number}</span>
            </li>
          ))}
        </ol>
      )}
    </section>
  );
}
