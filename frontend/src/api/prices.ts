export async function fetchPrices(): Promise<Prices> {
  const res = await fetch('/api/prices', { method: 'GET' });
  if (!res.ok) {
    throw new Error('No se pudo cargar los precios');
  }
  const prices: Prices = await res.json();
  return prices;
}

export async function updatePrices({
  montoBs,
  montoUsd,
}: Prices): Promise<void> {
  const res = await fetch('/api/prices', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ montoBs, montoUsd }),
  });
  if (!res.ok) {
    throw new Error('No se pudo actualizar los precios');
  }
}
