import { Link, useLocation } from 'react-router-dom';
import { FaShoppingCart, FaHome } from 'react-icons/fa';
import { AiOutlineLogout, AiOutlineClose } from 'react-icons/ai';
import { useAuth } from '@src/context/useAuth';

type SidebarProps = {
  open: boolean;
  setOpen: (o: boolean) => void;
};

export const Sidebar: React.FC<SidebarProps> = ({ open, setOpen }) => {
  const { logout } = useAuth();
  const location = useLocation();

  const menuItems = [
    {
      label: 'Compras',
      to: '/dashboard',
      icon: <FaShoppingCart />,
      isActive: location.pathname.startsWith('/dashboard'),
      tabIndex: 1,
    },
    {
      label: 'Pagina Principal',
      to: '/',
      icon: <FaHome />,
      isActive: location.pathname === '/',
      tabIndex: 2,
    },
  ];

  return (
    <>
      {/* Universal Overlay */}
      <div
        className={`
          fixed inset-0 z-[99] bg-black bg-opacity-50 transition-opacity
          ${
            open
              ? 'opacity-100 pointer-events-auto'
              : 'opacity-0 pointer-events-none'
          }
        `}
        onClick={() => setOpen(false)}
        aria-hidden={!open}
      />

      {/* Sidebar */}
      <aside
        className={`
          fixed top-0 left-0 z-[100] h-full w-60 bg-gray-900 text-white
          p-4 flex flex-col justify-between transform transition-transform duration-200
          ${open ? 'translate-x-0' : '-translate-x-full'}
        `}
        style={{ boxShadow: open ? '0 0 10px rgba(0,0,0,0.3)' : undefined }}
      >
        {/* Close button always */}
        <button
          className='mb-6 flex items-center'
          onClick={() => setOpen(false)}
          aria-label='Cerrar menú'
          type='button'
        >
          <AiOutlineClose className='h-6 w-6 text-white' />
          <span className='ml-2'>Cerrar menú</span>
        </button>

        <nav className='flex-1 flex flex-col gap-2 overflow-y-auto'>
          {menuItems.map((item) => (
            <Link
              key={item.label}
              to={item.to}
              tabIndex={item.tabIndex}
              className={`flex items-center gap-3 px-4 py-2 rounded-lg transition-colors ${
                item.isActive
                  ? 'bg-gray-800 font-bold'
                  : 'hover:bg-gray-800 font-normal'
              }`}
              onClick={() => setOpen(false)}
            >
              {item.icon}
              {item.label}
            </Link>
          ))}
        </nav>
        <div className='absolute bottom-4 left-0 w-full px-4'>
          <button
            tabIndex={3}
            onClick={() => logout.mutate()}
            className='flex items-center gap-3 w-full px-4 py-2 rounded-lg bg-red-700 hover:bg-red-800 font-bold transition-colors'
          >
            <AiOutlineLogout />
            Cerrar Sesion
          </button>
        </div>
      </aside>
    </>
  );
};
