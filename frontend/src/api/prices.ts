export async function fetchPrices(): Promise<Prices> {
  const res = await fetch('/api/prices', { method: 'GET' });
  if (!res.ok) {
    throw new Error('No se pudo cargar los precios');
  }
  const prices: Prices = await res.json();
  return prices;
}
