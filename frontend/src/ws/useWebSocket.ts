import { useEffect, useRef, useCallback } from 'react';
import { getSessionID } from '../api/http';
import { useGameStore } from '../store/gameStore';
import type { WSMessage, OutgoingMessage } from '../types';

export function useWebSocket(roomID: string | undefined) {
  const wsRef = useRef<WebSocket | null>(null);
  const setState = useGameStore((s) => s.setRoomState);
  const setError = useGameStore((s) => s.setError);

  useEffect(() => {
    if (!roomID) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_WS_URL || `${protocol}//${window.location.host}`;
    const url = `${host}/ws/${roomID}?session_id=${getSessionID()}`;

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      const msg: WSMessage = JSON.parse(event.data);
      if (msg.type === 'room_state' && msg.state) {
        setState(msg.state);
      } else if (msg.type === 'error' && msg.error) {
        setError(msg.error);
      }
    };

    ws.onclose = () => {
      wsRef.current = null;
    };

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [roomID, setState, setError]);

  const send = useCallback((msg: OutgoingMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    }
  }, []);

  return { send };
}
