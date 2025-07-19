import { useState, useEffect, type ReactNode } from 'react';
import { AuthContext } from './useAuth';
import { login as apiLogin } from '@src/api/auth';
import { getErrorMessage } from '@src/utils/errors';
export interface AuthContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  loading: boolean;
  authenticated: boolean;
  login: (email: string, password: string) => Promise<User>;
  logout: () => void;
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    fetch('/api/me', { credentials: 'include' })
      .then(async (res) => {
        if (res.status === 200) {
          const data = await res.json();
          setUser(data);
          setAuthenticated(true);
        } else {
          setUser(null);
          setAuthenticated(false);
        }
        setLoading(false);
      })
      .catch(() => {
        setUser(null);
        setAuthenticated(false);
        setLoading(false);
      });
  }, []);

  const login = async (email: string, password: string): Promise<User> => {
    setLoading(true);
    try {
      await apiLogin({ email, password });
      const userRes = await fetch('/api/me', { credentials: 'include' });
      if (userRes.status === 200) {
        const userData = await userRes.json();
        setUser(userData);
        setAuthenticated(true);
        setLoading(false);
        return userData;
      } else {
        setUser(null);
        setAuthenticated(false);
        setLoading(false);
        throw new Error('failed to retrieve user');
      }
    } catch (error) {
      setUser(null);
      setAuthenticated(false);
      setLoading(false);
      getErrorMessage(error, 'invalid login');
      throw error;
    }
  };

  const logout = () => {
    fetch('/api/logout', { method: 'POST', credentials: 'include' }).finally(
      () => {
        setUser(null);
        setAuthenticated(false);
      },
    );
  };

  return (
    <AuthContext.Provider
      value={{ user, setUser, loading, authenticated, login, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
};
