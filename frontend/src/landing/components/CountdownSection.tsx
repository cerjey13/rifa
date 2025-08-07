import { ProgressBar } from '@src/landing/components/ProgressBar';

type CountdownSectionProps = {
  percentage: number;
  loading: boolean;
  error: boolean;
};

const CountdownSection = ({
  percentage,
  loading,
  error,
}: CountdownSectionProps) => {
  if (loading)
    return <p className='text-gray-400'>Cargando disponibilidad...</p>;

  return (
    <div className='p-4 space-y-4 max-w-lg mx-auto'>
      <div className='flex flex-wrap items-center gap-3 text-brandLightGray'>
        <span className='flex items-center gap-1'>
          🕑La rifa se llevará a cabo luego de tener el 80% de los boletos
          vendidos.
        </span>
      </div>

      <ProgressBar percentage={percentage} />
      <p className='text-right text-brandLightGray'>
        {error ? 0 : percentage.toFixed(2)}% Vendido
      </p>
    </div>
  );
};

export default CountdownSection;
