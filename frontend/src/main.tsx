import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'sonner';
import { AppRouter } from '@src/AppRouter';
import { AuthProvider } from '@src/context/AuthProvider';
import { ModalProvider } from '@src/context/ModalProvider';
import '@src/index.css';

const queryClient = new QueryClient();

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <Toaster richColors position='top-right' />
      <div className='min-h-screen min-w-full bg-gray-900 text-white'>
        <AuthProvider>
          <ModalProvider>
            <AppRouter />
          </ModalProvider>
        </AuthProvider>
      </div>
    </QueryClientProvider>
  </React.StrictMode>,
);
