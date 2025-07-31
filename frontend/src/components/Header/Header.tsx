import { LoginRegisterModal } from '@src/components/AuthModal/Modal';
import { useAuth } from '@src/context/useAuth';
import { useModal } from '@src/context/useModal';

export const Header = () => {
  const { user, logout } = useAuth();
  const { modalOpen, isRegister, openLoginModal, closeModal } = useModal();

  const login = () => {
    closeModal();
  };

  return (
    <>
      <header className='sticky top-0 z-50 w-full flex justify-between items-center px-4 py-3 bg-[#121726cc] backdrop-blur-md shadow-md max-w-screen-xl mx-auto text-white'>
        <a href='/'>
          <img src='/logo.jpg' alt='Logo' className='h-10' />
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
            href='https://wa.me/+584141551801'
            target='_blank'
            rel='noreferrer'
            className='hover:text-[#F97316] transition-colors'
          >
            Whatsapp
          </a>

          {user && (
            <>
              {user.role === 'admin' && (
                <a href='/dashboard' className='hover:text-orange-400'>
                  Dashboard
                </a>
              )}
              <button
                onClick={() => logout.mutate()}
                className='hover:text-orange-400'
              >
                <a
                  href='/'
                  target='_self'
                  rel='noreferrer'
                  className='hover:text-[#F97316] transition-colors'
                >
                  Cerrar Sesión
                </a>
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
