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
