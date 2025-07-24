import { ProgressBar } from '@src/landing/components/ProgressBar';

export const CountdownSection = () => {
  //TODO: change the percentage from the backend
  const salesPercentage = 77.95;

  return (
    <div className='p-4 space-y-4 max-w-lg mx-auto'>
      <div className='flex flex-wrap items-center gap-3 text-brandLightGray'>
        <span className='flex items-center gap-1'>
          🕑La rifa se llevará a cabo luego de tener el 80% de los boletos
          vendidos.
        </span>
      </div>

      <ProgressBar percentage={salesPercentage} />
      <p className='text-right text-brandLightGray'>
        {salesPercentage}% vendido
      </p>
    </div>
  );
};
