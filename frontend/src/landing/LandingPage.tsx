import { PrizeInfo } from '@src/landing/components/PrizeInfo';
import { CountdownSection } from '@src/landing/components/CountdownSection';
import { Footer } from '@src/components/Footer/Footer';
import { Actions } from '@src/landing/components/Actions';
import { Header } from '@src/components/Header/Header';
import { useQuery } from '@tanstack/react-query';
import { fetchTicketsPercentage } from '@src/api/tickets';

export const LandingPage = () => {
  const {
    data: vendidos,
    isLoading,
    isError,
  } = useQuery<number, Error>({
    queryKey: ['ticket-availability'],
    queryFn: fetchTicketsPercentage,
    staleTime: 10 * 60 * 1000,
  });

  const safePercentage = isLoading || isError ? 0 : vendidos ?? 0;

  return (
    <>
      <Header />
      <main className='max-w-3xl mx-auto p-12 bg-gray-900 text-white min-h-auto'>
        <PrizeInfo />
        <CountdownSection
          percentage={safePercentage}
          loading={isLoading}
          error={isError}
        />
        <Actions />
      </main>
      <Footer />
    </>
  );
};
