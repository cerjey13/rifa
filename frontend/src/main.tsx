import React from 'react';
import ReactDOM from 'react-dom/client';
import { AppRouter } from './AppRouter';
import './index.css';
import { AuthProvider } from './context/AuthProvider';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <div className='min-h-screen min-w-full bg-gray-900 text-white'>
      <AuthProvider>
        <AppRouter />
      </AuthProvider>
    </div>
  </React.StrictMode>,
);
