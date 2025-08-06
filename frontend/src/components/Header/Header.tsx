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
        <div className='flex items-center gap-3'>
          <a href='/'>
            <img
              src='/logo.svg'
              alt='Logo'
              className='h-10 aspect-[1.22] rounded'
            />
          </a>
          {user && (
            <span
              title={user.email}
              className='text-sm sm:text-base font-medium text-white truncate max-w-[120px] sm:max-w-none'
            >
              Hola {user.name}!
            </span>
          )}
        </div>

        <nav className='flex items-center space-x-6 font-semibold text-sm'>
          {!user && (
            <button
              onClick={() => openLoginModal(false)}
              className='bg-[#C2410C] text-white px-4 py-2 rounded-md hover:bg-[#9A3412] transition'
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
                onClick={() => {
                  logout.mutate(undefined, {
                    onSuccess: () => {
                      setTimeout(() => window.location.reload(), 50);
                    },
                  });
                }}
                className='hover:text-orange-400'
              >
                Cerrar Sesión
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
