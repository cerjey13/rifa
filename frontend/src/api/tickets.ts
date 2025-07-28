export async function fetchUserTicketCount(): Promise<number> {
  const res = await fetch('/api/purchases/tickets', {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Error al cargar el conteo de tickets');
  const data = await res.json();
  return data.quantity;
}
