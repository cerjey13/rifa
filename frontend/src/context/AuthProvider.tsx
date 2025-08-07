import { useEffect, useState, type ReactNode } from 'react';
import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationResult,
} from '@tanstack/react-query';
import { login, fetchCurrentUser, logout } from '@src/api/auth';
import { AuthContext } from '@src/context/useAuth';

export interface AuthContextType {
  user: User | null;
  authenticated: boolean;
  loading: boolean;
  authReady: boolean;
  login: UseMutationResult<User, Error, { email: string; password: string }>;
  logout: UseMutationResult<void, Error, void>;
  refetchUser: () => void;
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [enableUserQuery, setEnableUserQuery] = useState(false);
  const [authReady, setAuthReady] = useState(false);

  useEffect(() => {
    requestIdleCallback(() => setEnableUserQuery(true));
  }, []);

  const queryClient = useQueryClient();

  const {
    data: user,
    isPending: userPending,
    isFetching: userFetching,
    error: userError,
    refetch: refetchUser,
  } = useQuery<User, Error>({
    queryKey: ['me'],
    queryFn: fetchCurrentUser,
    enabled: enableUserQuery,
    gcTime: 0,
    refetchOnWindowFocus: false,
    retry: false,
  });

  useEffect(() => {
    if (!userFetching && enableUserQuery) {
      setAuthReady(true);
    }
  }, [userFetching, enableUserQuery]);

  const loginMutation = useMutation<
    User,
    Error,
    { email: string; password: string }
  >({
    mutationFn: login,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['me'] });
    },
  });

  const logoutMutation = useMutation<void, Error, void>({
    mutationFn: logout,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['me'] });
    },
  });

  const authenticated = !!user && !userError;
  const loading =
    userPending ||
    userFetching ||
    loginMutation.isPending ||
    logoutMutation.isPending;

  return (
    <AuthContext.Provider
      value={{
        user: user || null,
        authenticated,
        loading,
        authReady,
        login: loginMutation,
        logout: logoutMutation,
        refetchUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
