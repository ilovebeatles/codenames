interface Props {
  clue: string;
  number: number;
  guessesLeft: number;
}

export default function ClueDisplay({ clue, number, guessesLeft }: Props) {
  if (!clue) return null;

  return (
    <div className="clue-display">
      <span className="clue-word">{clue}</span>
      <span className="clue-number">{number}</span>
      <span className="guesses-left">(осталось попыток: {guessesLeft})</span>
    </div>
  );
}
