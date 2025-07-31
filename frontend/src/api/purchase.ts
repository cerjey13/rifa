import { AuthError } from './auth';

export async function submitPurchase(purchase: {
  quantity: number;
  montoBs: string;
  montoUSD: string;
  paymentMethod: string;
  transactionDigits: string;
  selectedNumbers: string[];
  paymentScreenshot: File;
}): Promise<void> {
  const formData = new FormData();
  formData.append('quantity', String(purchase.quantity));
  formData.append('montoBs', purchase.montoBs);
  formData.append('montoUSD', purchase.montoUSD);
  formData.append('paymentMethod', purchase.paymentMethod);
  formData.append('transactionDigits', purchase.transactionDigits);
  formData.append(
    'selectedNumbers',
    purchase.selectedNumbers.filter(Boolean).length > 0
      ? purchase.selectedNumbers.filter(Boolean).join(',')
      : '',
  );
  formData.append('paymentScreenshot', purchase.paymentScreenshot);

  const res = await fetch('/api/purchases', {
    method: 'POST',
    body: formData,
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(await res.text());
  }
}

type PageParams = { page: number; perPage: number };
type PurchasesParams = PageParams & { status: string };

export async function fetchPurchases(
  params: PurchasesParams,
): Promise<{ purchase: Purchase[]; total: number }> {
  const urlParams = new URLSearchParams();
  if (params.status && params.status !== 'all') {
    urlParams.append('status', params.status);
  }
  urlParams.append('page', String(params.page));
  urlParams.append('perPage', String(params.perPage));

  const res = await fetch(`/api/purchases?${urlParams.toString()}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  if (res.status === 401) {
    throw new AuthError('Unauthorized');
  }
  if (!res.ok) throw new Error('Failed to fetch purchases');
  const purchase: Purchase[] = await res.json();
  const totalStr = res.headers.get('X-Total-Count') || '0';
  const total = parseInt(totalStr, 10);
  return { purchase, total };
}

export async function updatePurchaseStatus({
  transactionID,
  status,
}: {
  transactionID: string;
  status: PurchaseStatus;
}): Promise<void> {
  const res = await fetch(`/api/purchases?id=${transactionID}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ status }),
  });
  if (!res.ok) throw new Error('Network error');
}

export async function fetchMostPurchases(
  params: PageParams,
): Promise<{ purchases: MostPurchases[]; total: number }> {
  const urlParams = new URLSearchParams();
  urlParams.append('page', String(params.page));
  urlParams.append('perPage', String(params.perPage));

  const res = await fetch(
    `/api/purchases/leaderboard?${urlParams.toString()}`,
    {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    },
  );
  if (!res.ok) throw new Error('Error fetching leaderboard');
  const purchase: MostPurchases[] = await res.json();
  const totalStr = res.headers.get('X-Total-Count') || '0';
  const total = parseInt(totalStr, 10);
  return { purchases: purchase, total };
}
