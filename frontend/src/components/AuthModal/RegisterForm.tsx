import { register } from '@src/api/auth';
import { getErrorMessage } from '@src/utils/errors';
import { useState } from 'react';
import { FiEye, FiEyeOff } from 'react-icons/fi';

interface RegisterFormProps {
  onSwitch: () => void;
}

const validateEmail = (email: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

export const RegisterForm = ({ onSwitch }: RegisterFormProps) => {
  const [formData, setFormData] = useState({
    email: '',
    name: '',
    phone: '',
    password: '',
    passwordConfirm: '',
  });
  const [errors, setErrors] = useState({
    email: '',
    name: '',
    phone: '',

    password: '',
    passwordConfirm: '',
  });
  const [submitted, setSubmitted] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [errorMsg, setErrorMsg] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const validate = () => {
    const newErrors = {
      email: '',
      name: '',
      phone: '',
      password: '',
      passwordConfirm: '',
    };
    let valid = true;

    if (!validateEmail(formData.email)) {
      newErrors.email = 'Correo electrónico inválido';
      valid = false;
    }
    if (!formData.name.trim()) {
      newErrors.name = 'El nombre es obligatorio';
      valid = false;
    }
    if (!formData.phone.trim()) {
      newErrors.phone = 'El numero de telefono es obligatorio';
      valid = false;
    }
    if (formData.password.length < 6) {
      newErrors.password = 'La contraseña debe tener al menos 6 caracteres';
      valid = false;
    }
    if (formData.passwordConfirm !== formData.password) {
      newErrors.passwordConfirm = 'Las contraseñas no coinciden';
      valid = false;
    }

    setErrors(newErrors);
    return valid;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const togglePasswordVisibility = () => {
    setShowPassword((v) => !v);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);
    setSuccessMsg('');
    setErrorMsg('');

    if (!validate()) return;

    try {
      await register({
        name: formData.name,
        email: formData.email,
        password: formData.password,
        phone: formData.phone,
      });
      setSuccessMsg('¡Registro exitoso! Ya puedes iniciar sesión.');
      setTimeout(() => {
        onSwitch();
      }, 1000);
    } catch (error) {
      setErrorMsg(getErrorMessage(error, 'Error al registrarse'));
    }
  };

  return (
    <form onSubmit={handleSubmit} className='space-y-6'>
      <h2 className='text-2xl font-bold text-brandOrange'>Registrarse</h2>

      {successMsg && (
        <p className='text-green-500 font-semibold'>{successMsg}</p>
      )}
      {errorMsg && <p className='text-red-500 font-semibold'>{errorMsg}</p>}

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
        <label htmlFor='nombre' className='block mb-1 font-semibold text-white'>
          Nombre completo
        </label>
        <input
          id='nombre'
          name='name'
          type='text'
          placeholder='Tu nombre completo'
          value={formData.name}
          onChange={handleChange}
          className={`w-full bg-[#1E2638] border rounded px-3 py-2 text-white placeholder-brandLightGray focus:outline-none focus:ring-2 focus:ring-brandOrange ${
            errors.name && submitted
              ? 'border-red-500'
              : 'border-brandLightGray'
          }`}
        />
        {errors.name && submitted && (
          <p className='text-red-500 mt-1 text-sm'>{errors.name}</p>
        )}
      </div>

      <div>
        <label htmlFor='phone' className='block mb-1 font-semibold text-white'>
          Telefono
        </label>
        <input
          id='phone'
          name='phone'
          type='number'
          placeholder='Tu numero de telefono'
          value={formData.phone}
          onChange={handleChange}
          className={`w-full bg-[#1E2638] border rounded px-3 py-2 text-white placeholder-brandLightGray focus:outline-none focus:ring-2 focus:ring-brandOrange ${
            errors.phone && submitted
              ? 'border-red-500'
              : 'border-brandLightGray'
          }`}
        />
        {errors.phone && submitted && (
          <p className='text-red-500 mt-1 text-sm'>{errors.phone}</p>
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
            placeholder='Mínimo 8 caracteres'
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

      <div>
        <label
          htmlFor='passwordConfirm'
          className='block mb-1 font-semibold text-white'
        >
          Repetir contraseña
        </label>
        <input
          id='passwordConfirm'
          name='passwordConfirm'
          type='password'
          placeholder='Repite tu contraseña'
          value={formData.passwordConfirm}
          onChange={handleChange}
          className={`w-full bg-[#1E2638] border rounded px-3 py-2 text-white placeholder-brandLightGray focus:outline-none focus:ring-2 focus:ring-brandOrange ${
            errors.passwordConfirm && submitted
              ? 'border-red-500'
              : 'border-brandLightGray'
          }`}
        />
        {errors.passwordConfirm && submitted && (
          <p className='text-red-500 mt-1 text-sm'>{errors.passwordConfirm}</p>
        )}
      </div>

      <button
        type='submit'
        className='w-full bg-[#FF7F00] text-white py-2 rounded hover:bg-orange-600 transition'
      >
        Registrarse
      </button>

      <p className='text-center text-sm text-brandLightGray'>
        ¿Ya tienes cuenta?{' '}
        <button
          type='button'
          onClick={onSwitch}
          className='text-[#FF7F00] hover:underline'
        >
          Inicia Sesión
        </button>
      </p>
    </form>
  );
};
