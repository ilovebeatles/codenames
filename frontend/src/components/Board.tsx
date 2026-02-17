import type { CardView } from '../types';
import Card from './Card';

interface Props {
  cards: CardView[];
  onGuess: (cardID: string) => void;
  canGuess: boolean;
}

export default function Board({ cards, onGuess, canGuess }: Props) {
  const sorted = [...cards].sort((a, b) => a.position - b.position);

  return (
    <div className="board">
      {sorted.map((card) => (
        <Card
          key={card.id}
          card={card}
          onClick={() => onGuess(card.id)}
          disabled={!canGuess}
        />
      ))}
    </div>
  );
}
