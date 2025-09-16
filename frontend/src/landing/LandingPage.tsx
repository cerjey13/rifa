import { PrizeInfo } from '@src/landing/components/PrizeInfo';
import { Footer } from '@src/components/Footer/Footer';
import { Header } from '@src/components/Header/Header';
import { useQueries } from '@tanstack/react-query';
import { fetchTicketsPercentage } from '@src/api/tickets';
import React, { Suspense } from 'react';
import { fetchPrices } from '@src/api/prices';
import premio512 from '@src/assets/premio-500.webp';

const CountdownSection = React.lazy(
  () => import('@src/landing/components/CountdownSection'),
);

const Actions = React.lazy(() => import('@src/landing/components/Actions'));

export const LandingPage = () => {
  const results = useQueries({
    queries: [
      {
        queryKey: ['ticket-availability'],
        queryFn: fetchTicketsPercentage,
        staleTime: 10 * 60 * 1000,
      },
      {
        queryKey: ['prices'],
        queryFn: fetchPrices,
        staleTime: 10 * 60 * 1000,
      },
    ],
  });

  const [vendidos, prices] = results;

  const safePercentage =
    vendidos.isLoading || vendidos.isError ? 0 : vendidos.data ?? 0;
  const safePrices =
    prices.isLoading || prices.isError ? undefined : prices.data ?? undefined;

  return (
    <>
      {/* preload the pize image */}
      <link
        rel='preload'
        as='image'
        href={premio512}
        type='image/webp'
        imageSrcSet={`${premio512} 384w`}
        imageSizes='(max-width: 384px) 100vw, 384px'
        fetchPriority='high'
      />
      <Header />
      <main className='max-w-3xl mx-auto p-12 bg-gray-900 text-white min-h-auto'>
        <PrizeInfo />
        <Suspense fallback={null}>
          <CountdownSection
            percentage={safePercentage}
            loading={vendidos.isLoading}
            error={vendidos.isError}
          />
        </Suspense>
        <Suspense fallback={null}>
          <Actions prices={safePrices} />
        </Suspense>
      </main>
      <Footer />
    </>
  );
};
