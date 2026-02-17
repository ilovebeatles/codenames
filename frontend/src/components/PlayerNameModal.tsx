import { useState } from 'react';

interface Props {
  onSubmit: (name: string) => void;
}

export default function PlayerNameModal({ onSubmit }: Props) {
  const [name, setName] = useState('');

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Введите ваше имя</h2>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Имя"
          maxLength={20}
          autoFocus
          onKeyDown={(e) => {
            if (e.key === 'Enter' && name.trim()) onSubmit(name.trim());
          }}
        />
        <button
          onClick={() => name.trim() && onSubmit(name.trim())}
          disabled={!name.trim()}
        >
          Войти
        </button>
      </div>
    </div>
  );
}
