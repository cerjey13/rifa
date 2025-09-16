import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AdminApp } from '@src/admin/Admin';
import { LandingPage } from '@src/landing/LandingPage';
import { AuthGuard } from '@src/context/AuthGuard';

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path='/' element={<LandingPage />} />
        <Route
          path='/dashboard/*'
          element={
            <AuthGuard>
              <AdminApp />
            </AuthGuard>
          }
        />
        <Route path='*' element={<Navigate to='/' replace />} />
      </Routes>
    </BrowserRouter>
  );
};
