import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { fetchSearchByNumber } from '@src/api/tickets';

export const Searcher: React.FC = () => {
  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');
  const [selectedUser, setSelectedUser] = useState<Omit<User, 'role'> | null>(
    null,
  );

  const { data: result, isFetching } = useQuery({
    queryKey: ['search-number', submittedQuery],
    queryFn: () => fetchSearchByNumber(submittedQuery),
    enabled: !!submittedQuery,
  });

  const isValidQuery = /^\d{1,4}$/.test(query) && +query >= 0 && +query <= 9999;
  const handleSearch = () => {
    if (isValidQuery) {
      setSubmittedQuery(query.trim());
    }
  };

  return (
    <div className='overflow-x-auto bg-gray-900 rounded-lg shadow p-4 space-y-4'>
      <div className='flex gap-2 items-center'>
        <input
          type='text'
          inputMode='numeric'
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder='Buscar número (0 - 9999)...'
          className='px-3 py-2 rounded bg-gray-800 text-white w-full md:max-w-sm'
        />
        <button
          onClick={handleSearch}
          disabled={!isValidQuery || isFetching}
          className='px-4 py-2 bg-orange-500 hover:bg-orange-600 text-white font-semibold rounded disabled:opacity-50'
        >
          {isFetching ? 'Buscando...' : 'Buscar'}
        </button>
      </div>

      {!isValidQuery && query !== '' && (
        <p className='text-sm text-red-400'>
          Ingresa un número válido entre 0 y 9999.
        </p>
      )}

      <table className='min-w-full divide-y divide-gray-700'>
        <thead className='bg-gray-800'>
          <tr>
            <th className='px-2 md:px-4 py-2 text-left font-medium'>Usuario</th>
            <th className='px-2 md:px-4 py-2 text-left font-medium'>
              Números comprados
            </th>
          </tr>
        </thead>
        <tbody className='divide-y divide-gray-700'>
          {result && (
            <tr key={result.user.id}>
              <td className='px-2 md:px-4 py-2'>
                <button
                  className='text-orange-400 hover:underline'
                  onClick={() => setSelectedUser(result.user)}
                >
                  {result.user.name}
                </button>
              </td>
              <td className='px-2 md:px-4 py-2'>{result.tickets.join(', ')}</td>
            </tr>
          )}
          {!result && !isFetching && submittedQuery && (
            <tr>
              <td colSpan={2} className='text-center px-4 py-6 text-gray-400'>
                No hay resultados para "{submittedQuery}".
              </td>
            </tr>
          )}
        </tbody>
      </table>

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
