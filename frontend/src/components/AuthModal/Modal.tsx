import { LoginForm } from '@src/components/AuthModal/LoginForm';
import { RegisterForm } from '@src/components/AuthModal/RegisterForm';

interface LoginRegisterModalProps {
  isRegister: boolean;
  onClose: () => void;
  onSwitch: () => void;
  onLogin: (user: User) => void;
}

export const LoginRegisterModal = ({
  isRegister,
  onClose,
  onSwitch,
  onLogin,
}: LoginRegisterModalProps) => {
  return (
    <div className='fixed inset-0 bg-black bg-opacity-70 flex justify-center items-center z-60'>
      <div className='bg-[#121726] rounded-md shadow-lg w-96 p-6 relative text-white'>
        <button
          onClick={onClose}
          className='absolute top-3 right-3 text-brandLightGray hover:text-white font-bold text-xl'
          aria-label='Close modal'
        >
          &times;
        </button>

        {isRegister ? (
          <RegisterForm onLogin={onLogin} />
        ) : (
          <LoginForm onLogin={onLogin} onSwitch={onSwitch} />
        )}
      </div>
    </div>
  );
};
