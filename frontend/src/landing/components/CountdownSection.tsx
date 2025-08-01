import { fetchTicketsPercentage } from '@src/api/tickets';
import { ProgressBar } from '@src/landing/components/ProgressBar';
import { useQuery } from '@tanstack/react-query';

export const CountdownSection = () => {
  const {
    data: vendidos,
    isLoading,
    isError,
  } = useQuery<number, Error>({
    queryKey: ['ticket-availability'],
    queryFn: () => fetchTicketsPercentage(),
    staleTime: 10 * 60 * 1000,
  });

  if (isLoading)
    return <p className='text-gray-400'>Cargando disponibilidad...</p>;
  // if (isError || !vendidos)
  //  {vendidos = 0}

  return (
    <div className='p-4 space-y-4 max-w-lg mx-auto'>
      <div className='flex flex-wrap items-center gap-3 text-brandLightGray'>
        <span className='flex items-center gap-1'>
          🕑La rifa se llevará a cabo luego de tener el 80% de los boletos
          vendidos.
        </span>
      </div>

      <ProgressBar percentage={vendidos ?? 0} />
      <p className='text-right text-brandLightGray'>{vendidos}% Vendido</p>
    </div>
  );
};
