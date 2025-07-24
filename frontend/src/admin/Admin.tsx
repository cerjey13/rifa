import { useState } from 'react';
import { Sidebar } from '@src/admin/components/SideBar';
import { FaBars } from 'react-icons/fa';
import { Dashboard } from '@src/admin/components/Dashboard';

export const AdminApp: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className='flex min-h-screen bg-gray-900 text-white'>
      {/* Sidebar */}
      <Sidebar open={sidebarOpen} setOpen={setSidebarOpen} />

      {/* Main content */}
      <main className='flex-1 p-2 sm:p-4 md:p-6 bg-gray-900'>
        {/* Mobile: Hamburger button + header */}
        <div className='flex items-center gap-2 mb-4 md:mb-6'>
          <button
            className='md:hidden flex-shrink-0'
            onClick={() => setSidebarOpen(true)}
            aria-label='Abrir menú'
            type='button'
          >
            <FaBars className='h-7 w-7 text-white' />
          </button>
          <h1 className='text-xl sm:text-2xl md:text-3xl font-bold'>
            Resumen de compras
          </h1>
        </div>
        <Dashboard />
      </main>
    </div>
  );
};
