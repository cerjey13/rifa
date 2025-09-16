import { useAuth } from './useAuth';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export const AuthGuard = ({ children }: { children: React.ReactNode }) => {
  const { authenticated, loading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && !authenticated) {
      navigate('/', { replace: true });
    }
  }, [authenticated, loading, navigate]);

  // Optionally show a loader
  if (loading) return <div className="text-center mt-10">Cargando...</div>;

  return <>{children}</>;
};
