import type { CardView } from '../types';

interface Props {
  card: CardView;
  onClick: () => void;
  disabled: boolean;
}

const TYPE_COLORS: Record<string, string> = {
  red: '#e74c3c',
  blue: '#3498db',
  neutral: '#bdc3c7',
  assassin: '#2c3e50',
};

export default function Card({ card, onClick, disabled }: Props) {
  const opacity = card.card_type === 'assassin' ? '99' : '30'
  const bgColor = card.revealed
    ? TYPE_COLORS[card.card_type] || '#f5f5f5'
    : card.card_type && !card.revealed
      ? `${TYPE_COLORS[card.card_type]}${opacity}` // spymaster hint
      : '#f5f5f5';

  const textColor = card.revealed && (card.card_type === 'assassin' || card.card_type === 'red' || card.card_type === 'blue')
    ? '#fff'
    : '#333';

  return (
    <button
      className={`card ${card.revealed ? 'revealed' : ''} ${card.card_type && !card.revealed ? 'hinted' : ''}`}
      style={{ backgroundColor: bgColor, color: textColor }}
      onClick={onClick}
      disabled={disabled || card.revealed}
    >
      {card.word}
    </button>
  );
}
