export async function submitPurchase(purchase: {
  quantity: number;
  montoBs: string;
  montoUSD: string;
  paymentMethod: string;
  transactionDigits: string;
  paymentScreenshot: File;
}): Promise<void> {
  const formData = new FormData();
  formData.append('quantity', String(purchase.quantity));
  formData.append('montoBs', purchase.montoBs);
  formData.append('montoUSD', purchase.montoUSD);
  formData.append('paymentMethod', purchase.paymentMethod);
  formData.append('transactionDigits', purchase.transactionDigits);
  formData.append('paymentScreenshot', purchase.paymentScreenshot);

  const res = await fetch('/api/purchase', {
    method: 'POST',
    body: formData,
    credentials: 'include',
  });
  if (!res.ok) {
    throw new Error(await res.text());
  }
}

type PageParams = { page: number; perPage: number; status: string };

export async function fetchPurchases(params: PageParams): Promise<Purchase[]> {
  const urlParams = new URLSearchParams();
  if (params.status && params.status !== 'all')
    urlParams.append('status', params.status);
  urlParams.append('page', String(params.page));
  urlParams.append('perPage', String(params.perPage));

  const res = await fetch(`/api/purchases?${urlParams.toString()}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Network error');
  return res.json();
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
