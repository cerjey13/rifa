import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import AdminApp from '@src/admin/App';
import { LandingPage } from '@src/landing/LandingPage';

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path='/' element={<LandingPage />} />
        <Route path='/dashboard/*' element={<AdminApp />} />
        <Route path='*' element={<Navigate to='/' replace />} />
      </Routes>
    </BrowserRouter>
  );
};
