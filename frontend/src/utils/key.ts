export async function buildIdempotencyKey(purchase: {
  userId: string;
  quantity: number;
  montoBs: string;
  montoUSD: string;
  paymentMethod: string;
  transactionDigits: string;
  selectedNumbers: string[];
}): Promise<string> {
  const payload = JSON.stringify({
    userId: purchase.userId,
    quantity: purchase.quantity,
    montoBs: purchase.montoBs,
    montoUSD: purchase.montoUSD,
    paymentMethod: purchase.paymentMethod,
    transactionDigits: purchase.transactionDigits,
    selectedNumbers: purchase.selectedNumbers.sort(),
  });

  const buffer = new TextEncoder().encode(payload);
  const digest = await crypto.subtle.digest('SHA-256', buffer);
  const hashArray = Array.from(new Uint8Array(digest));
  return hashArray.map((b) => b.toString(16).padStart(2, '0')).join('');
}
