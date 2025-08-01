interface User {
  id: string;
  email: string;
  name: string;
  phone: string;
  role: 'admin' | 'user';
}

type PurchaseStatus = 'pending' | 'verified' | 'cancelled';

interface Purchase {
  user: Omit<User, 'role'>;
  id: string;
  quantity: number;
  tickets: string[];
  montoBs: number;
  montoUsd: number;
  paymentMethod: string;
  transactionDigits: string;
  selectedNumbers: string[];
  paymentScreenshot: string;
  status: string;
  date: string;
}

type MostPurchases = {
  user: Omit<User, 'role'>;
  quantity: number;
};
