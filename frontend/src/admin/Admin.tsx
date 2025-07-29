import { useState } from 'react';
import { Sidebar } from '@src/admin/components/SideBar';
import { FaBars } from 'react-icons/fa';
import { Dashboard } from '@src/admin/components/Dashboard';

export const AdminApp: React.FC = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className='relative flex min-h-screen bg-gray-900 text-white'>
      <Sidebar open={sidebarOpen} setOpen={setSidebarOpen} />

      <button
        className='fixed z-[200] top-4 left-4 md:top-6 md:left-6 flex-shrink-0 bg-transparent'
        onClick={() => setSidebarOpen(true)}
        aria-label='Abrir menú'
        type='button'
        style={{
          display: sidebarOpen ? 'none' : 'flex',
        }}
      >
        <FaBars className='h-7 w-7 text-white drop-shadow' />
      </button>

      <main className='flex-1 p-4 md:p-6 bg-gray-900'>
        <div className='flex items-center gap-2 mb-4 pl-14'>
          <h1 className='text-2xl md:text-2xl font-bold'>Resumen de compras</h1>
        </div>
        <Dashboard />
      </main>
    </div>
  );
};
