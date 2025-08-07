import React, { useEffect, useState } from 'react';
import { fetchPrices, updatePrices } from '@src/api/prices';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

export const Precios: React.FC = () => {
  const [bs, setBs] = useState<string>('');
  const [usd, setUsd] = useState<string>('');
  const [modalOpen, setModalOpen] = useState<boolean>(false);

  const queryClient = useQueryClient();
  const {
    data: prices,
    isLoading,
    isError,
    error,
  } = useQuery<Prices, Error>({
    queryKey: ['prices'],
    queryFn: fetchPrices,
    staleTime: 10 * 60 * 1000,
  });

  useEffect(() => {
    if (prices) {
      setBs(prices.montoBs.toString());
      setUsd(prices.montoUsd.toString());
    }
  }, [prices]);

  const mutation = useMutation({
    mutationFn: updatePrices,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['prices'] });
    },
  });

  const handleConfirmClick = () => {
    setModalOpen(true);
  };

  // If user confirms in the modal
  const handleModalConfirm = () => {
    setModalOpen(false);
    mutation.mutate({
      montoBs: Number(bs),
      montoUsd: Number(usd),
    });
  };

  const handleModalCancel = () => {
    setModalOpen(false);
  };

  return (
    <div className='p-2 sm:p-4 md:p-6'>
      <div className='bg-gray-800 rounded-xl p-8 max-w-md mx-auto text-center shadow-lg'>
        <h2 className='text-xl font-bold mb-4'>Gestión de Precios</h2>
        <p className='text-gray-300 mb-6'>
          Modificar los precios por ticket de la rifa:
        </p>
        {isLoading ? (
          <div className='py-8 text-gray-400'>Cargando precios...</div>
        ) : isError ? (
          <div className='py-8 text-red-400'>
            Error al cargar precios: {error?.message}
          </div>
        ) : (
          <>
            <div className='flex flex-col gap-4 mb-6'>
              <label className='flex flex-col text-left'>
                <span className='mb-1'>Precio en Bolívares (Bs):</span>
                <input
                  type='number'
                  value={bs}
                  min={0}
                  step='0.01'
                  onChange={(e) => setBs(e.target.value)}
                  className='rounded-lg px-4 py-2 bg-gray-900 border border-gray-700 focus:outline-none focus:ring focus:border-blue-500 transition'
                  disabled={mutation.isPending}
                />
              </label>
              <label className='flex flex-col text-left'>
                <span className='mb-1'>Precio en Dólares (USD):</span>
                <input
                  type='number'
                  value={usd}
                  min={0}
                  step='0.01'
                  onChange={(e) => setUsd(e.target.value)}
                  className='rounded-lg px-4 py-2 bg-gray-900 border border-gray-700 focus:outline-none focus:ring focus:border-blue-500 transition'
                  disabled={mutation.isPending}
                />
              </label>
            </div>
            <button
              onClick={handleConfirmClick}
              disabled={mutation.isPending || !bs || !usd}
              className='bg-blue-600 hover:bg-blue-700 rounded-xl px-8 py-2 font-bold text-white transition disabled:opacity-50'
            >
              {mutation.isPending ? 'Guardando...' : 'Confirmar Cambios'}
            </button>
            {mutation.isSuccess && (
              <div className='mt-4 text-green-400'>
                Precios actualizados exitosamente.
              </div>
            )}
            {mutation.isError && (
              <div className='mt-4 text-red-400'>
                Error al actualizar precios: {(mutation.error as Error).message}
              </div>
            )}
          </>
        )}
        {modalOpen && (
          <div className='fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-60'>
            <div className='bg-gray-800 rounded-xl shadow-xl p-6 max-w-xs w-full text-center'>
              <h3 className='text-lg font-bold mb-4'>¿Confirmar cambios?</h3>
              <p className='mb-6 text-gray-300'>
                Vas a cambiar el precio a:
                <br />
                <span className='font-semibold'>
                  Bs {bs} / USD {usd}
                </span>
              </p>
              <div className='flex justify-center gap-4'>
                <button
                  onClick={handleModalConfirm}
                  className='bg-green-600 hover:bg-green-700 px-6 py-2 rounded-lg font-bold'
                  disabled={mutation.isPending}
                >
                  Confirmar
                </button>
                <button
                  onClick={handleModalCancel}
                  className='bg-gray-600 hover:bg-gray-700 px-6 py-2 rounded-lg font-bold'
                  disabled={mutation.isPending}
                >
                  Cancelar
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
