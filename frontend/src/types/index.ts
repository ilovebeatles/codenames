export type Team = 'red' | 'blue' | '';
export type Role = 'spymaster' | 'operative' | '';
export type Phase = 'lobby' | 'playing' | 'finished';

export interface Room {
  id: string;
  created_at: string;
}

export interface Player {
  id: string;
  room_id: string;
  name: string;
  team: Team;
  role: Role;
  is_online: boolean;
}

export interface Game {
  id: string;
  room_id: string;
  phase: Phase;
  current_team: Team;
  current_clue: string;
  current_number: number;
  guesses_left: number;
  winner: Team;
}

export interface CardView {
  id: string;
  word: string;
  card_type: string; // red | blue | neutral | assassin | '' (hidden)
  position: number;
  revealed: boolean;
  revealed_by: Team;
}

export interface RoomState {
  room: Room;
  players: Player[];
  game: Game | null;
  cards: CardView[];
}

export interface WSMessage {
  type: string;
  state?: RoomState;
  error?: string;
}

export interface OutgoingMessage {
  type: string;
  team?: string;
  role?: string;
  clue?: string;
  number?: number;
  card_id?: string;
}
