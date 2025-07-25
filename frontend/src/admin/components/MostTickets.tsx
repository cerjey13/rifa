import { fetchMostPurchases } from '@src/api/purchase';
import { safeArray } from '@src/utils/arrays';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';

export const TicketsLeaderboard: React.FC = () => {
  const [selectedUser, setSelectedUser] = useState<Omit<User, 'role'> | null>(
    null,
  );
  const [page, setPage] = useState<number>(1);
  const perPage = 10;

  const {
    data: mostPurchasesRaw,
    isLoading,
    isFetching,
    error,
  } = useQuery<{ purchases: MostPurchases[]; total: number }, Error>({
    queryKey: ['most-purchases', { page, perPage }],
    queryFn: () => fetchMostPurchases({ page, perPage }),
    placeholderData: (prevData) => prevData,
  });

  const mostPurchases = safeArray<MostPurchases>(mostPurchasesRaw?.purchases);
  const pageCount = Math.max(
    1,
    Math.ceil((mostPurchasesRaw?.total ?? 0) / perPage),
  );

  if (isLoading) return <div className='text-center mt-10'>Cargando...</div>;
  if (error)
    return (
      <div className='text-center mt-10 text-red-500'>Error cargando datos</div>
    );

  return (
    <div className='overflow-x-auto bg-gray-900 rounded-lg shadow'>
      {isFetching && (
        <div className='flex items-center gap-2 text-sm text-gray-400 mb-2'>
          <svg
            className='animate-spin h-4 w-4 text-gray-400'
            fill='none'
            viewBox='0 0 24 24'
          >
            <circle
              className='opacity-25'
              cx='12'
              cy='12'
              r='10'
              stroke='currentColor'
              strokeWidth='4'
            ></circle>
            <path
              className='opacity-75'
              fill='currentColor'
              d='M4 12a8 8 0 018-8v8H4z'
            ></path>
          </svg>
          Actualizando datos...
        </div>
      )}
      <table className='min-w-full divide-y divide-gray-700'>
        <thead className='bg-gray-800'>
          <tr>
            <th className='px-2 md:px-4 py-2 text-left font-medium'>Usuario</th>
            <th className='px-2 md:px-4 py-2 text-left font-medium'>
              Tickets comprados
            </th>
          </tr>
        </thead>
        <tbody className='divide-y divide-gray-700'>
          {mostPurchases.map((p) => (
            <tr key={p.user.id}>
              <td className='px-2 md:px-4 py-2 flex items-center gap-2'>
                {p.user.name}
                <button
                  className='ml-2 px-2 py-1 rounded bg-gray-700 text-xs text-white hover:bg-gray-600'
                  onClick={() => setSelectedUser(p.user)}
                  title='Ver detalles'
                  type='button'
                >
                  Detalles
                </button>
              </td>
              <td className='px-2 md:px-4 py-2 font-semibold'>{p.quantity}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className='flex gap-2 my-4'>
        <button
          disabled={page === 1}
          onClick={() => setPage((p) => Math.max(p - 1, 1))}
          className='px-3 py-1 rounded bg-gray-700 disabled:bg-gray-800'
        >
          Anterior
        </button>
        <span>
          Página {page} de {pageCount}
        </span>
        <button
          disabled={page === pageCount}
          onClick={() => setPage((p) => Math.min(p + 1, pageCount))}
          className='px-3 py-1 rounded bg-gray-700 disabled:bg-gray-800'
        >
          Siguiente
        </button>
      </div>
      {selectedUser && (
        <div
          className='fixed inset-0 z-[10000] flex items-center justify-center bg-black bg-opacity-80'
          onClick={() => setSelectedUser(null)}
        >
          <div
            className='relative bg-gray-900 p-6 rounded-xl shadow-xl min-w-[320px]'
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className='text-lg font-bold mb-4'>Detalles del usuario</h2>
            <div className='space-y-2'>
              <div>
                <span className='font-semibold'>Nombre:</span>{' '}
                {selectedUser.name}
              </div>
              <div>
                <span className='font-semibold'>Email:</span>{' '}
                {selectedUser.email}
              </div>
              <div>
                <span className='font-semibold'>Teléfono:</span>{' '}
                {selectedUser.phone}
              </div>
            </div>
            <button
              className='absolute top-2 right-2 bg-white bg-opacity-90 rounded-full p-1 hover:bg-opacity-100 transition'
              onClick={() => setSelectedUser(null)}
              aria-label='Cerrar'
              type='button'
            >
              <svg
                className='h-5 w-5 text-black'
                fill='none'
                viewBox='0 0 24 24'
                stroke='currentColor'
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  strokeWidth={2}
                  d='M6 18L18 6M6 6l12 12'
                />
              </svg>
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
