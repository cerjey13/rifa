import { PrizeInfo } from '@src/landing/components/PrizeInfo';
import { CountdownSection } from '@src/landing/components/CountdownSection';
import { Footer } from '@src/components/Footer/Footer';
import { Actions } from '@src/landing/components/Actions';
import { Header } from '@src/components/Header/Header';

export const LandingPage = () => {
  return (
    <>
      <Header />
      <main className='max-w-3xl mx-auto p-12 bg-gray-900 text-white min-h-auto'>
        <PrizeInfo />
        <CountdownSection />
        <Actions />
      </main>
      <Footer />
    </>
  );
};
