import { useAuth } from '@src/context/useAuth';
import { getErrorMessage } from '@src/utils/errors';
import { useState } from 'react';
import { FiEye, FiEyeOff } from 'react-icons/fi';

interface LoginFormProps {
  onLogin: () => void;
  onSwitch: () => void;
}

interface LoginProps {
  email: string;
  password: string;
}

interface FormError {
  email: string;
  password: string;
  message: string;
}
const validarEmail = (email: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

export const LoginForm = ({ onLogin, onSwitch }: LoginFormProps) => {
  const [formData, setFormData] = useState<LoginProps>({
    email: '',
    password: '',
  });
  const [errors, setErrors] = useState<FormError>({
    email: '',
    password: '',
    message: '',
  });
  const [submitted, setSubmitted] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [showPassword, setShowPassword] = useState(false);
  const { login } = useAuth();

  const validar = () => {
    const newErrors: FormError = { email: '', password: '', message: '' };
    let valid = true;

    if (!validarEmail(formData.email)) {
      newErrors.email = 'Correo electrónico inválido';
      valid = false;
    }
    if (!formData.password) {
      newErrors.password = 'La contraseña es obligatoria';
      valid = false;
    }

    setErrors(newErrors);
    return valid;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const togglePasswordVisibility = () => {
    setShowPassword((v) => !v);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);
    setErrors((prev) => ({ ...prev, message: '' }));

    if (!validar()) return;

    setLoading(true);

    try {
      const user = await login.mutateAsync({
        email: formData.email,
        password: formData.password,
      });

      if (user.email !== undefined) {
        onLogin();
      } else {
        setErrors((prev) => ({ ...prev, message: 'Error al iniciar sesión' }));
      }
    } catch (error) {
      console.error(getErrorMessage(error, 'Error al iniciar sesión'));
      setErrors((prev) => ({
        ...prev,
        message: 'Error al iniciar sesión',
      }));
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className='space-y-6'>
      <h2 className='text-2xl font-bold text-brandOrange'>Iniciar sesión</h2>

      {errors.message && (
        <p className='text-red-500 text-center font-semibold'>
          {errors.message}
        </p>
      )}

      <div>
        <label htmlFor='email' className='block mb-1 font-semibold text-white'>
          Correo electrónico
        </label>
        <input
          id='email'
          name='email'
          type='email'
          placeholder='tu@email.com'
          value={formData.email}
          onChange={handleChange}
          className={`w-full bg-[#1E2638] border rounded px-3 py-2 text-white placeholder-brandLightGray focus:outline-none focus:ring-2 focus:ring-brandOrange ${
            errors.email && submitted
              ? 'border-red-500'
              : 'border-brandLightGray'
          }`}
        />
        {errors.email && submitted && (
          <p className='text-red-500 mt-1 text-sm'>{errors.email}</p>
        )}
      </div>

      <div>
        <label
          htmlFor='password'
          className='block mb-1 font-semibold text-white'
        >
          Contraseña
        </label>
        <div className='relative'>
          <input
            id='password'
            name='password'
            type={showPassword ? 'text' : 'password'}
            placeholder='Tu contraseña'
            value={formData.password}
            onChange={handleChange}
            className={`w-full bg-[#1E2638] border rounded px-3 py-2 text-white placeholder-brandLightGray focus:outline-none focus:ring-2 focus:ring-brandOrange ${
              errors.password && submitted
                ? 'border-red-500'
                : 'border-brandLightGray'
            }`}
          />
          <button
            type='button'
            className='absolute right-2 top-1/2 -translate-y-1/2 text-lg text-white'
            tabIndex={-1}
            onClick={togglePasswordVisibility}
            aria-label={
              showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'
            }
          >
            {showPassword ? <FiEyeOff /> : <FiEye />}
          </button>
        </div>
        {errors.password && submitted && (
          <p className='text-red-500 mt-1 text-sm'>{errors.password}</p>
        )}
      </div>

      <button
        type='submit'
        disabled={loading}
        className='w-full bg-[#FF7F00] text-white py-2 rounded hover:bg-orange-600 transition disabled:opacity-50 disabled:cursor-not-allowed'
      >
        {loading ? 'Cargando...' : 'Entrar'}
      </button>

      <p className='text-center text-sm text-brandLightGray'>
        ¿No tienes cuenta?{' '}
        <button
          type='button'
          onClick={onSwitch}
          className='text-[#FF7F00] hover:underline'
        >
          Regístrate
        </button>
      </p>
    </form>
  );
};
