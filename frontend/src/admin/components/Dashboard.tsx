import { useEffect, useState } from 'react';

export const Dashboard = () => {
  const [purchases, setPurchases] = useState<Purchase[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalImage, setModalImage] = useState<string | null>(null);

  useEffect(() => {
    fetch('/api/purchases', {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    })
      .then((res) => res.json())
      .then((data) => {
        setPurchases(data);
        setLoading(false);
      });
  }, []);

  if (loading) return <div className='text-center mt-10'>Cargando...</div>;

  const totalBs = purchases.reduce((sum, p) => sum + (p.montoBs || 0), 0);
  const totalUsd = purchases.reduce((sum, p) => sum + (p.montoUsd || 0), 0);

  return (
    <div className='p-2 sm:p-4 md:p-6'>
      <h1 className='text-xl sm:text-2xl md:text-3xl font-bold mb-4 md:mb-6'>
        Resumen de compras
      </h1>

      {/* Stats */}
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

      {/* Responsive Table */}
      <div className='block md:hidden'>
        {/* Mobile: show as cards */}
        <div className='space-y-4'>
          {purchases.map((p) => (
            <div
              key={p.userId + p.transactionDigits}
              className='bg-gray-900 rounded-lg shadow p-4'
            >
              <div className='flex flex-wrap items-center gap-2 mb-2'>
                <span className='font-bold'>Usuario:</span> {p.userId}
              </div>
              <div>
                Cantidad: <span className='font-semibold'>{p.quantity}</span>
              </div>
              <div>
                Monto Bs: <span className='font-semibold'>{p.montoBs}</span>
              </div>
              <div>
                Monto USD: <span className='font-semibold'>{p.montoUsd}</span>
              </div>
              <div>Método: {p.paymentMethod}</div>
              <div>Transacción: {p.transactionDigits}</div>
              <div>
                Estado: <span className='capitalize'>{p.status}</span>
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
              <th className='px-2 md:px-4 py-2 text-left font-medium'>
                Captura
              </th>
            </tr>
          </thead>
          <tbody className='divide-y divide-gray-700'>
            {purchases.map((p) => (
              <tr key={p.userId + p.transactionDigits}>
                <td className='px-2 md:px-4 py-2'>{p.userId}</td>
                <td className='px-2 md:px-4 py-2'>{p.quantity}</td>
                <td className='px-2 md:px-4 py-2'>{p.montoBs}</td>
                <td className='px-2 md:px-4 py-2'>{p.montoUsd}</td>
                <td className='px-2 md:px-4 py-2'>{p.paymentMethod}</td>
                <td className='px-2 md:px-4 py-2'>{p.transactionDigits}</td>
                <td className='px-2 md:px-4 py-2 capitalize'>{p.status}</td>
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
    </div>
  );
};
