import type { Dispatch, SetStateAction } from 'react';

interface PaymentMethodsProps {
  selectedMethod: string | null;
  setSelectedMethod: Dispatch<SetStateAction<string | null>>;
  onBack: () => void;
  onClose: () => void;
  onNext: () => void;
}

const paymentMethods = {
  pagoMovil: 'pago movil',
  zelle: 'zelle',
} as const;

const paymentOptions = [
  { id: paymentMethods.pagoMovil, label: 'PAGOMOVIL', icon: '📱' },
  {
    id: paymentMethods.zelle,
    label: paymentMethods.zelle.toUpperCase(),
    icon: '💸',
  },
];

export const PaymentMethods = ({
  selectedMethod,
  setSelectedMethod,
  onBack,
  onClose,
  onNext,
}: PaymentMethodsProps) => {
  const handleClose = () => {
    setSelectedMethod(null);
    onClose();
  };
  return (
    <div className='flex flex-col gap-6  max-h-[90vh] overflow-y-auto py-4'>
      <div className='flex items-center justify-between'>
        <h2 className='font-bold uppercase text-lg'>Métodos de pago</h2>

        <button
          onClick={handleClose}
          aria-label='Cerrar'
          className='p-2 text-white hover:text-yellow-400 transition'
        >
          <svg
            xmlns='http://www.w3.org/2000/svg'
            className='h-6 w-6'
            fill='none'
            viewBox='0 0 24 24'
            stroke='currentColor'
            strokeWidth={2}
          >
            <path
              strokeLinecap='round'
              strokeLinejoin='round'
              d='M6 18L18 6M6 6l12 12'
            />
          </svg>
        </button>
      </div>
      <hr className='border-gray-700' />

      <div className='grid grid-cols-1 sm:grid-cols-3 gap-4'>
        {paymentOptions.map(({ id, label, icon }) => (
          <button
            key={id}
            onClick={() =>
              selectedMethod === id
                ? setSelectedMethod(null)
                : setSelectedMethod(id)
            }
            className={`flex flex-col items-center p-4 rounded bg-gray-800 hover:bg-gray-700 transition ${
              selectedMethod === id ? 'ring-2 ring-yellow-500' : ''
            }`}
          >
            <span className='text-4xl'>{icon}</span>
            <span className='mt-2 font-semibold'>{label}</span>
            <span className='mt-2 bg-gray-700 rounded px-2 py-1 text-sm uppercase'>
              Selecciona aquí
            </span>
          </button>
        ))}
      </div>

      <div className='flex gap-4 justify-end'>
        <button
          onClick={onBack}
          className='bg-gray-700 py-2 px-6 rounded hover:bg-gray-600 transition'
        >
          Atrás
        </button>
        <button
          onClick={onNext}
          disabled={!selectedMethod}
          className='bg-yellow-500 py-2 px-6 rounded hover:bg-yellow-600 text-black font-bold transition disabled:opacity-50 disabled:cursor-not-allowed'
        >
          Siguiente
        </button>
      </div>
    </div>
  );
};
