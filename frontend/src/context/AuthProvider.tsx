import type { ReactNode } from 'react';
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
  login: UseMutationResult<User, Error, { email: string; password: string }>;
  logout: UseMutationResult<void, Error, void>;
  refetchUser: () => void;
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const queryClient = useQueryClient();

  const {
    data: user,
    isLoading: userLoading,
    error: userError,
    refetch: refetchUser,
  } = useQuery<User, Error>({
    queryKey: ['me'],
    queryFn: fetchCurrentUser,
    refetchOnWindowFocus: false,
    retry: false,
  });

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
    userLoading || loginMutation.isPending || logoutMutation.isPending;

  return (
    <AuthContext.Provider
      value={{
        user: user || null,
        authenticated,
        loading,
        login: loginMutation,
        logout: logoutMutation,
        refetchUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
