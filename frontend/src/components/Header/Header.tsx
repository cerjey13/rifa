import logo from '@src/assets/react.svg';
import { LoginRegisterModal } from '@src/components/AuthModal/Modal';
import { useAuth } from '@src/context/useAuth';
import { useModal } from '@src/context/useModal';

export const Header = () => {
  const { user, setUser, logout } = useAuth();
  const { modalOpen, isRegister, openLoginModal, closeModal } = useModal();

  const login = (user: User) => {
    setUser(user);
    closeModal();
  };

  return (
    <>
      <header className='sticky top-0 z-50 w-full flex justify-between items-center px-4 py-3 bg-[#121726cc] backdrop-blur-md shadow-md max-w-screen-xl mx-auto text-white'>
        <a href='/'>
          <img src={logo} alt='Logo' className='h-10' />
        </a>
        <nav className='flex items-center space-x-6 font-semibold text-sm'>
          {!user && (
            <button
              onClick={() => openLoginModal(false)}
              className='bg-[#F97316] text-white px-4 py-2 rounded-md hover:bg-[#EA580C] transition'
            >
              Iniciar sesión / Registrarse
            </button>
          )}

          <a
            href='https://wa.me/your-number'
            target='_blank'
            rel='noreferrer'
            className='hover:text-[#F97316] transition-colors'
          >
            Whatsapp
          </a>
          <a
            href='/events'
            className='hover:text-brandOrange transition-colors'
          >
            Eventos
          </a>

          {user && (
            <>
              {user.role === 'admin' && (
                <a href='/dashboard' className='hover:text-brandOrange'>
                  Dashboard
                </a>
              )}
              <button onClick={logout} className='hover:text-brandOrange'>
                Cerrar sesión
              </button>
            </>
          )}
        </nav>
      </header>

      {modalOpen && (
        <LoginRegisterModal
          isRegister={isRegister}
          onClose={closeModal}
          onSwitch={() => openLoginModal(!isRegister)}
          onLogin={login}
        />
      )}
    </>
  );
};
