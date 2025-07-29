interface TicketResponse {
  tickets: string[];
}

export async function fetchUserTickets(): Promise<string[]> {
  const res = await fetch('/api/purchases/tickets', {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Error al cargar el conteo de tickets');
  const data: TicketResponse = await res.json();
  return data.tickets;
}
