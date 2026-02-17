import { useState } from 'react';

interface Props {
  onGiveClue: (clue: string, number: number) => void;
}

export default function ClueInput({ onGiveClue }: Props) {
  const [clue, setClue] = useState('');
  const [number, setNumber] = useState(1);

  const handleSubmit = () => {
    if (clue.trim()) {
      onGiveClue(clue.trim(), number);
      setClue('');
      setNumber(1);
    }
  };

  return (
    <div className="clue-input">
      <input
        type="text"
        value={clue}
        onChange={(e) => setClue(e.target.value)}
        placeholder="Подсказка..."
        onKeyDown={(e) => {
          if (e.key === 'Enter') handleSubmit();
        }}
      />
      <select value={number} onChange={(e) => setNumber(Number(e.target.value))}>
        {[0, 1, 2, 3, 4, 5, 6, 7, 8, 9].map((n) => (
          <option key={n} value={n}>{n}</option>
        ))}
      </select>
      <button onClick={handleSubmit} disabled={!clue.trim()}>
        Дать подсказку
      </button>
    </div>
  );
}
