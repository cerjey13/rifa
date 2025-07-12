import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';

import { LandingPage } from './landing/LandingPage';
import App from './admin/App';

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path='/' element={<LandingPage />} />
        <Route path='/admin/*' element={<App />} />
        <Route path='*' element={<Navigate to='/' replace />} />
      </Routes>
    </BrowserRouter>
  );
};
