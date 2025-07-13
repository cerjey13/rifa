import { useState } from 'react';

interface RegisterFormProps {
  onLogin: (user: User) => void;
}

const validarEmail = (email: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

export const RegisterForm = ({ onLogin }: RegisterFormProps) => {
  const [formData, setFormData] = useState({
    email: '',
    name: '',
    password: '',
    passwordConfirm: '',
  });
  const [errors, setErrors] = useState({
    email: '',
    name: '',
    password: '',
    passwordConfirm: '',
  });
  const [submitted, setSubmitted] = useState(false);

  const validar = () => {
    const newErrors = {
      email: '',
      name: '',
      password: '',
      passwordConfirm: '',
    };
    let valid = true;

    if (!validarEmail(formData.email)) {
      newErrors.email = 'Correo electrónico inválido';
      valid = false;
    }
    if (!formData.name.trim()) {
      newErrors.name = 'El nombre es obligatorio';
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
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);
    if (!validar()) return;

    onLogin({
      name: formData.name,
      email: formData.email,
      role: 'user',
      phone: '555',
    });
  };

  return (
    <form onSubmit={handleSubmit} className='space-y-6'>
      <h2 className='text-2xl font-bold text-brandOrange'>Registrarse</h2>

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
          name='nombre'
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
        <label
          htmlFor='password'
          className='block mb-1 font-semibold text-white'
        >
          Contraseña
        </label>
        <input
          id='password'
          name='password'
          type='password'
          placeholder='Mínimo 6 caracteres'
          value={formData.password}
          onChange={handleChange}
          className={`w-full bg-[#1E2638] border rounded px-3 py-2 text-white placeholder-brandLightGray focus:outline-none focus:ring-2 focus:ring-brandOrange ${
            errors.password && submitted
              ? 'border-red-500'
              : 'border-brandLightGray'
          }`}
        />
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
    </form>
  );
};
