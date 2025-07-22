interface User {
  email: string;
  name: string;
  phone: string;
  role: 'admin' | 'user';
}

interface Purchase {
  user: Omit<User, 'role'>;
  quantity: number;
  montoBs: number;
  montoUsd: number;
  paymentMethod: string;
  transactionDigits: string;
  paymentScreenshot: string;
  status: string;
  date: string;
}
