const API_BASE = import.meta.env.VITE_API_URL || '';

function getSessionID(): string {
  let sid = localStorage.getItem('session_id');
  if (!sid) {
    sid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (crypto.getRandomValues(new Uint8Array(1))[0] & 15);
      const v = c === 'x' ? r : (r & 0x3) | 0x8;
      return v.toString(16);
    });
    localStorage.setItem('session_id', sid);
  }
  return sid;
}

export { getSessionID };

export async function createRoom(): Promise<{ id: string }> {
  const res = await fetch(`${API_BASE}/api/rooms`, {
    method: 'POST',
    headers: { 'X-Session-ID': getSessionID() },
  });
  if (!res.ok) throw new Error('Failed to create room');
  return res.json();
}

export async function joinRoom(roomID: string, name: string): Promise<void> {
  const res = await fetch(`${API_BASE}/api/players`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Session-ID': getSessionID(),
    },
    body: JSON.stringify({ room_id: roomID, name }),
  });
  if (!res.ok) throw new Error('Failed to join room');
}

export async function getRoomState(roomID: string) {
  const res = await fetch(`${API_BASE}/api/rooms/${roomID}`, {
    headers: { 'X-Session-ID': getSessionID() },
  });
  if (!res.ok) throw new Error('Room not found');
  return res.json();
}
