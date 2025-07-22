import React from 'react';
import ReactDOM from 'react-dom/client';
import { AppRouter } from '@src/AppRouter';
import { AuthProvider } from '@src/context/AuthProvider';
import { ModalProvider } from '@src/context/ModalProvider';
import '@src/index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <div className='min-h-screen min-w-full bg-gray-900 text-white'>
      <AuthProvider>
        <ModalProvider>
          <AppRouter />
        </ModalProvider>
      </AuthProvider>
    </div>
  </React.StrictMode>,
);
