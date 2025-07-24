export async function fetchUserTicketCount(email: string): Promise<number> {
  const params = new URLSearchParams({ email });
  const res = await fetch(`/api/user-ticket-count?${params.toString()}`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Error al cargar el conteo de tickets');
  const data = await res.json();
  return data.count;
}
