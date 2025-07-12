import { type ReactNode, useState } from 'react';
import { AuthContext } from './useAuth';

export interface AuthContextType {
  user: User | null;
  setUser: (user: User | null) => void;
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);

  return (
    <AuthContext.Provider value={{ user, setUser }}>
      {children}
    </AuthContext.Provider>
  );
};
