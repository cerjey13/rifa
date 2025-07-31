interface TicketResponse {
  tickets: string[];
}

export async function fetchUserTickets(): Promise<string[]> {
  const res = await fetch('/api/tickets/users', {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Error al cargar el conteo de tickets');
  const data: TicketResponse = await res.json();
  return data.tickets;
}

interface SearchResult {
  user: Omit<User, 'role'>;
  tickets: string[];
}

export async function fetchSearchByNumber(
  num: string,
): Promise<SearchResult | null> {
  const urlParams = new URLSearchParams();
  urlParams.append('number', num);
  const res = await fetch(`/api/purchases/search?${urlParams.toString()}`, {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Error buscando el número');
  return await res.json();
}

export async function fetchAvailableTickets(num: string[]): Promise<string[]> {
  const urlParams = new URLSearchParams();
  urlParams.append('numbers', num.join(','));

  const res = await fetch(`/api/tickets?${urlParams.toString()}`, {
    method: 'GET',
    credentials: 'include',
  });

  if (!res.ok) throw new Error('Error buscando números');
  const data = await res.json();
  return data.tickets;
}

export async function fetchTicketsPercentage(): Promise<number> {
  const res = await fetch('/api/tickets/percentage', { method: 'GET' });
  if (!res.ok) {
    throw new Error('No se pudo cargar el porcentaje de disponibilidad');
  }
  const data = await res.json();
  return data.vendidos;
}
