import { useState, useEffect } from 'react';

interface TicketsModalProps {
  userEmail: string;
  onClose: () => void;
}

export const TicketsModal = ({ userEmail, onClose }: TicketsModalProps) => {
  const [loading, setLoading] = useState(false);
  const [ticketsCount, setTicketsCount] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTicketsCount = async () => {
      setLoading(true);
      setError(null);

      try {
        // const response = await fetch(
        //   `/api/tickets?email=${encodeURIComponent(userEmail)}`,
        // );
        // if (!response.ok) throw new Error('Error fetching tickets');

        // const data = await response.json();
        setTicketsCount(100);
      } catch (error) {
        setError('No se pudo obtener la cantidad de números comprados.');
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    fetchTicketsCount();
  }, [userEmail]);

  return (
    <div className='fixed inset-0 flex items-center justify-center bg-black bg-opacity-70 z-50 p-4'>
      <div className='bg-gray-800 p-6 rounded-md text-white max-w-sm w-full relative'>
        <div className='flex items-center justify-between mb-4'>
          <h2 className='font-bold uppercase text-lg'>Mis números comprados</h2>

          <button
            onClick={onClose}
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
        <hr className='border-gray-700 mb-4' />

        {loading && (
          <p className='mb-4 text-center text-yellow-400 font-semibold'>
            Cargando...
          </p>
        )}

        {error && (
          <p className='mb-4 text-red-500 text-center font-semibold'>{error}</p>
        )}

        {!loading && ticketsCount !== null && (
          <p className='text-center text-lg'>
            Has comprado <strong>{ticketsCount}</strong> número(s).
          </p>
        )}
      </div>
    </div>
  );
};
