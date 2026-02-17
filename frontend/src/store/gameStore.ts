import { create } from 'zustand';
import type { RoomState, Player } from '../types';

interface GameStore {
  roomState: RoomState | null;
  error: string | null;
  playerName: string | null;
  setRoomState: (state: RoomState) => void;
  setError: (error: string) => void;
  clearError: () => void;
  setPlayerName: (name: string) => void;
  currentPlayer: () => Player | null;
}

export const useGameStore = create<GameStore>((set, get) => ({
  roomState: null,
  error: null,
  playerName: localStorage.getItem('player_name'),

  setRoomState: (state) => set({ roomState: state, error: null }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: null }),
  setPlayerName: (name) => {
    localStorage.setItem('player_name', name);
    set({ playerName: name });
  },

  currentPlayer: () => {
    const state = get().roomState;
    if (!state) return null;
    return state.players.find((p) => p.name === get().playerName) ?? null;
  },
}));
