import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Admin, Resource } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';

const dataProvider = simpleRestProvider(import.meta.env.VITE_API_URL || '');

const AdminApp = () => (
  <Admin dataProvider={dataProvider}>
    <Resource name='tickets' />
    {/* other admin resources */}
  </Admin>
);

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path='/admin/*'
          element={
            // TODO: Add auth guard here to protect admin routes
            <AdminApp />
          }
        />
        <Route path='*' element={<Navigate to='/' replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
