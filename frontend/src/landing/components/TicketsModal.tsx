import { fetchUserTickets } from '@src/api/tickets';
import { useQuery } from '@tanstack/react-query';

interface TicketsModalProps {
  userEmail: string;
  onClose: () => void;
}

export const TicketsModal = ({ userEmail, onClose }: TicketsModalProps) => {
  const {
    data: tickets,
    isLoading,
    error,
  } = useQuery<string[], Error>({
    queryKey: ['ticket-count', userEmail],
    queryFn: () => fetchUserTickets(),
    enabled: !!userEmail,
  });

  if (!userEmail) return null;

  return (
    <div className='fixed inset-0 flex items-center justify-center bg-black bg-opacity-70 z-50 p-4'>
      <div className='bg-gray-800 p-6 rounded-md text-white max-w-xs w-full relative shadow-lg'>
        <div className='flex items-center justify-between mb-4'>
          <h2 className='font-bold uppercase text-base sm:text-lg'>
            Mis números comprados
          </h2>
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

        {isLoading ? (
          <p className='text-center text-gray-400 text-base sm:text-lg'>
            Cargando...
          </p>
        ) : error ? (
          <p className='text-center text-red-500 text-base sm:text-lg'>
            Error al cargar tickets
          </p>
        ) : tickets && tickets.length === 0 ? (
          <p className='text-center text-gray-400 text-base sm:text-lg'>
            No has comprado ningún número.
          </p>
        ) : (
          tickets && (
            <div className='space-y-4'>
              <p className='text-center text-base'>
                Has comprado{' '}
                <strong className='text-yellow-400 text-lg'>
                  {tickets.length}
                </strong>{' '}
                número(s):
              </p>
              <div className='bg-gray-700 p-2 rounded-lg shadow-inner max-h-[200px] overflow-y-auto'>
                <div className='flex flex-wrap gap-1.5 p-2'>
                  {tickets.map((ticket) => (
                    <span
                      key={ticket}
                      className='bg-yellow-400 text-black font-mono font-semibold w-13 text-center py-2 rounded-full text-sm shadow'
                    >
                      {ticket}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          )
        )}
      </div>
    </div>
  );
};
