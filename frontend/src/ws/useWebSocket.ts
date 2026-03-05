import { useEffect, useRef, useCallback } from 'react';
import { getSessionID } from '../api/http';
import { useGameStore } from '../store/gameStore';
import type { WSMessage, OutgoingMessage } from '../types';

const RECONNECT_DELAYS = [500, 1000, 2000, 5000];

export function useWebSocket(roomID: string | undefined) {
  const wsRef = useRef<WebSocket | null>(null);
  const roomIDRef = useRef<string | undefined>(undefined);
  const reconnectAttempt = useRef<number>(0);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const setState = useGameStore((s) => s.setRoomState);
  const setError = useGameStore((s) => s.setError);

  useEffect(() => {
    if (!roomID) return;

    // If already connected/connecting to the same room, skip (handles StrictMode re-mount)
    if (roomIDRef.current === roomID && wsRef.current && wsRef.current.readyState <= WebSocket.OPEN) {
      return;
    }

    // Clean up old connection if room changed
    clearTimeout(reconnectTimer.current);
    if (wsRef.current) {
      wsRef.current.onclose = null;
      wsRef.current.close();
      wsRef.current = null;
    }

    roomIDRef.current = roomID;

    function connect() {
      if (roomIDRef.current !== roomID) return;

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = import.meta.env.VITE_WS_URL || `${protocol}//${window.location.host}`;
      const url = `${host}/ws/${roomID}?session_id=${getSessionID()}`;

      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        reconnectAttempt.current = 0;
      };

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
        if (roomIDRef.current === roomID) {
          const delay = RECONNECT_DELAYS[Math.min(reconnectAttempt.current, RECONNECT_DELAYS.length - 1)];
          reconnectAttempt.current++;
          reconnectTimer.current = setTimeout(connect, delay);
        }
      };
    }

    connect();
  }, [roomID, setState, setError]);

  // Close WS on true unmount (component removed from tree)
  useEffect(() => {
    return () => {
      roomIDRef.current = undefined;
      clearTimeout(reconnectTimer.current);
      if (wsRef.current) {
        wsRef.current.onclose = null;
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, []);

  const send = useCallback((msg: OutgoingMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    }
  }, []);

  return { send };
}
