import { useEffect, useState } from 'react';
import { PurchaseFilters } from './PurchaseFilters';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { fetchPurchases, updatePurchaseStatus } from '@src/api/purchase';
import { toast } from 'sonner';
import { formatDateES } from '@src/utils/dates';
import { safeArray } from '@src/utils/arrays';
import { useNavigate } from 'react-router-dom';
import { AuthError } from '@src/api/auth';

type Filters = { status: string };
const statusEs: Record<string, string> = {
  pending: 'pendiente',
  verified: 'verificado',
  cancelled: 'cancelado',
};

export const ResumenCompras: React.FC = () => {
  const [modalImage, setModalImage] = useState<string | null>(null);
  const [selectedUser, setSelectedUser] = useState<Omit<User, 'role'> | null>(
    null,
  );
  const [filters, setFilters] = useState<Filters>({ status: '' });
  const [page, setPage] = useState<number>(1);
  const [editingStatus, setEditingStatus] = useState<
    Record<string, PurchaseStatus>
  >({});
  const perPage = 10;
  const navigate = useNavigate();

  const {
    data: purchasesRaw,
    isLoading,
    isFetching,
    error,
  } = useQuery<{ purchase: Purchase[]; total: number }, Error>({
    queryKey: ['purchases', { ...filters, page, perPage }],
    queryFn: () => fetchPurchases({ ...filters, page, perPage }),
    placeholderData: (prevData) => prevData,
    staleTime: 1000 * 60 * 5,
  });
  const queryClient = useQueryClient();

  const { mutate: changeStatus, isPending } = useMutation({
    mutationFn: updatePurchaseStatus,
    onSuccess: (_, variables) => {
      const label = statusEs[variables.status];
      const capitalized = label.charAt(0).toUpperCase() + label.slice(1);
      toast.success(`Estado actualizado a "${capitalized}"`);
      queryClient.invalidateQueries({ queryKey: ['purchases'] });
    },
    onError: (error: Error) => {
      console.error('Update failed:', error);
      toast.error('No se pudo actualizar el estado. Intenta nuevamente.');
    },
  });

  useEffect(() => {
    if (error && error instanceof AuthError) {
      navigate('/', { replace: true });
    }
  }, [error, navigate]);

  const purchases = safeArray<Purchase>(purchasesRaw?.purchase);
  const pageCount = Math.max(
    1,
    Math.ceil((purchasesRaw?.total ?? 0) / perPage),
  );

  if (isLoading) return <div className='text-center mt-10'>Cargando...</div>;
  if (error)
    return (
      <div className='text-center mt-10 text-red-500'>Error cargando datos</div>
    );

  const totalBs = purchases.reduce((sum, p) => sum + (p.montoBs || 0), 0) ?? 0;
  const totalUsd =
    purchases.reduce((sum, p) => sum + (p.montoUsd || 0), 0) ?? 0;

  return (
    <div className='p-2 sm:p-4 md:p-6'>
      {/* Stats - mobile stacks, desktop grid */}
      <div className='grid grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-6 mb-4 md:mb-8'>
        <div className='bg-gray-800 rounded-lg shadow p-4 flex flex-col items-center'>
          <span className='text-base sm:text-lg'>Total compras</span>
          <span className='text-xl sm:text-2xl font-bold'>
            {purchases.length}
          </span>
        </div>
        <div className='bg-gray-800 rounded-lg shadow p-4 flex flex-col items-center'>
          <span className='text-base sm:text-lg'>Total Bs</span>
          <span className='text-xl sm:text-2xl font-bold'>
            {totalBs.toLocaleString()}
          </span>
        </div>
        <div className='bg-gray-800 rounded-lg shadow p-4 flex flex-col items-center'>
          <span className='text-base sm:text-lg'>Total USD</span>
          <span className='text-xl sm:text-2xl font-bold'>
            ${totalUsd.toFixed(2)}
          </span>
        </div>
      </div>

      <PurchaseFilters onFilter={setFilters} />

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

      {/* MOBILE FIRST: Cards */}
      <div className='block md:hidden'>
        <div className='space-y-4'>
          {purchases.map((p) => (
            <div key={p.user.id} className='bg-gray-900 rounded-lg shadow p-4'>
              <div className='flex flex-wrap items-center gap-2 mb-2'>
                <span className='font-bold'>Usuario:</span> {p.user.name}
                <button
                  className='ml-2 px-2 py-1 rounded bg-gray-700 text-xs text-white hover:bg-gray-600'
                  onClick={() => setSelectedUser(p.user)}
                  title='Ver detalles'
                  type='button'
                >
                  Detalles
                </button>
              </div>
              <div>
                <span className='font-semibold'>Cantidad:</span> {p.quantity}
              </div>
              <div>
                <span className='font-semibold'>Monto Bs:</span> {p.montoBs}
              </div>
              <div>
                <span className='font-semibold'>Monto USD:</span> {p.montoUsd}
              </div>
              <div>
                <span className='font-semibold'>Método:</span> {p.paymentMethod}
              </div>
              <div>
                <span className='font-semibold'>Transacción:</span>{' '}
                {p.transactionDigits}
              </div>
              <div>
                <span className='font-semibold'>Estado:</span>{' '}
                <div className='flex flex-col gap-1 mt-1 w-fit'>
                  <select
                    value={editingStatus[p.id] ?? p.status}
                    onChange={(e) =>
                      setEditingStatus((prev) => ({
                        ...prev,
                        [p.id]: e.target.value as PurchaseStatus,
                      }))
                    }
                    disabled={isPending}
                    className='bg-gray-800 text-white rounded px-2 py-1 capitalize'
                  >
                    {Object.entries(statusEs).map(([value, label]) => (
                      <option value={value} key={value}>
                        {label}
                      </option>
                    ))}
                  </select>
                  {editingStatus[p.id] && editingStatus[p.id] !== p.status && (
                    <button
                      className='px-2 py-1 rounded bg-blue-600 text-xs text-white hover:bg-blue-700'
                      onClick={() => {
                        console.log(editingStatus[p.id]);
                        changeStatus({
                          transactionID: p.id,
                          status: editingStatus[p.id],
                        });
                        setEditingStatus((prev) => {
                          const copy = { ...prev };
                          delete copy[p.id];
                          return copy;
                        });
                      }}
                      disabled={isPending}
                      type='button'
                    >
                      Confirmar
                    </button>
                  )}
                </div>
              </div>
              <div>
                <span className='font-semibold'>Fecha:</span>{' '}
                {formatDateES(p.date)}
              </div>
              <div>
                {p.paymentScreenshot ? (
                  <img
                    src={`data:image/png;base64,${p.paymentScreenshot}`}
                    alt='Comprobante'
                    className='h-14 w-14 object-contain border rounded shadow mt-2 cursor-pointer'
                    onClick={() =>
                      setModalImage(
                        `data:image/png;base64,${p.paymentScreenshot}`,
                      )
                    }
                  />
                ) : (
                  <span className='text-gray-500'>Sin captura</span>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Desktop Table */}
      <div className='hidden md:block overflow-x-auto bg-gray-900 rounded-lg shadow'>
        <table className='min-w-full divide-y divide-gray-700'>
          <thead className='bg-gray-800'>
            <tr>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Usuario
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Cantidad
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Monto Bs
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Monto USD
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Método
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Transacción
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Estado
              </th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>Fecha</th>
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Captura
              </th>
            </tr>
          </thead>
          <tbody className='divide-y divide-gray-700'>
            {purchases.map((p) => (
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
                <td className='px-2 md:px-4 py-2'>{p.quantity}</td>
                <td className='px-2 md:px-4 py-2'>{p.montoBs}</td>
                <td className='px-2 md:px-4 py-2'>{p.montoUsd}</td>
                <td className='px-2 md:px-4 py-2'>{p.paymentMethod}</td>
                <td className='px-2 md:px-4 py-2'>{p.transactionDigits}</td>
                <td className='px-2 md:px-4 py-2'>
                  <select
                    value={editingStatus[p.id] ?? p.status}
                    onChange={(e) =>
                      setEditingStatus((prev) => ({
                        ...prev,
                        [p.id]: e.target.value as PurchaseStatus,
                      }))
                    }
                    disabled={isPending}
                    className='bg-gray-800 text-white capitalize rounded px-2 py-1'
                  >
                    {Object.entries(statusEs).map(([value, label]) => (
                      <option className='capitalize' value={value} key={value}>
                        {label}
                      </option>
                    ))}
                  </select>
                  {editingStatus[p.id] && editingStatus[p.id] !== p.status && (
                    <button
                      className='mt-1 px-2 py-1 rounded bg-blue-600 text-xs text-white hover:bg-blue-700'
                      onClick={() => {
                        console.log(editingStatus[p.id]);
                        changeStatus({
                          transactionID: p.id,
                          status: editingStatus[p.id],
                        });
                        setEditingStatus((prev) => {
                          const copy = { ...prev };
                          delete copy[p.id];
                          return copy;
                        });
                      }}
                      disabled={isPending}
                      type='button'
                    >
                      Confirmar
                    </button>
                  )}
                </td>
                <td className='px-2 md:px-4 py-2 capitalize'>
                  {formatDateES(p.date)}
                </td>
                <td className='px-2 md:px-4 py-2'>
                  {p.paymentScreenshot ? (
                    <img
                      src={`data:image/png;base64,${p.paymentScreenshot}`}
                      alt='Comprobante'
                      className='h-12 w-12 object-contain border rounded shadow cursor-pointer'
                      onClick={() =>
                        setModalImage(
                          `data:image/png;base64,${p.paymentScreenshot}`,
                        )
                      }
                    />
                  ) : (
                    <span className='text-gray-500'>—</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination Controls */}
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

      {/* Modal for large image */}
      {modalImage && (
        <div
          className='fixed inset-0 z-[9999] flex items-center justify-center bg-black bg-opacity-80'
          onClick={() => setModalImage(null)}
        >
          <div className='relative flex items-center justify-center w-full h-full'>
            <img
              src={modalImage}
              alt='Comprobante ampliado'
              className='max-w-[90vw] max-h-[80vh] rounded-xl shadow-xl'
              onClick={(e) => e.stopPropagation()}
            />
            <button
              className='absolute top-15 right-5 z-[10000] bg-white bg-opacity-90 rounded-full p-2 hover:bg-opacity-100 transition'
              onClick={() => setModalImage(null)}
              aria-label='Cerrar'
              type='button'
            >
              <svg
                xmlns='http://www.w3.org/2000/svg'
                className='h-7 w-7 text-black'
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

      {/* Modal for user details */}
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
