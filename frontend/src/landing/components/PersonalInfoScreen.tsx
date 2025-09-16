import { useState } from 'react';

interface PersonalInfoScreenProps {
  user: User | null;
  onBack: () => void;
  onClose: () => void;
  onNext: (data: { name: string; phone: string; email: string }) => void;
}

//TODO: see if we need to use this anywhere?
export const PersonalInfoScreen = ({
  user,
  onBack,
  onClose,
  onNext,
}: PersonalInfoScreenProps) => {
  const [name, setName] = useState(user?.name || '');
  const [phone, setPhone] = useState(user?.phone || '');
  const [email, setEmail] = useState(user?.email || '');

  const handleNext = () => {
    // Optionally validate fields here
    if (!name || !phone || !email) {
      alert('Por favor, completa todos los campos');
      return;
    }
    // You could add more validations (email format, phone format, etc.)
    onNext({ name, phone, email });
  };

  return (
    <div className='relative flex flex-col gap-6'>
      <button
        onClick={onClose}
        aria-label='Cerrar'
        className='absolute top-0 right-0 p-2 text-white hover:text-yellow-400 transition'
      >
        ✕
      </button>
      <h2 className='text-xl font-bold uppercase'>Información personal</h2>
      <hr className='border-gray-700' />

      <label className='flex flex-col'>
        Nombre y Apellido
        <input
          type='text'
          required
          value={name}
          onChange={(e) => setName(e.target.value)}
          className='bg-gray-700 rounded p-2 text-white focus:outline-none focus:ring-2 focus:ring-yellow-500'
        />
      </label>

      <label className='flex flex-col'>
        Número de teléfono
        <input
          type='tel'
          required
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          className='bg-gray-700 rounded p-2 text-white focus:outline-none focus:ring-2 focus:ring-yellow-500'
        />
      </label>

      <label className='flex flex-col'>
        Correo
        <input
          type='email'
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className='bg-gray-700 rounded p-2 text-white focus:outline-none focus:ring-2 focus:ring-yellow-500'
        />
      </label>

      <div className='flex gap-4 justify-end'>
        <button
          type='button'
          onClick={onBack}
          className='bg-gray-700 py-2 px-6 rounded hover:bg-gray-600 transition'
        >
          Atrás
        </button>
        <button
          type='button'
          onClick={handleNext}
          className='bg-yellow-500 py-2 px-6 rounded hover:bg-yellow-600 text-black font-bold transition'
        >
          Siguiente
        </button>
      </div>
    </div>
  );
};
