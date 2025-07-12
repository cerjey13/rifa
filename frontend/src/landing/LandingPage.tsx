import { PrizeInfo } from '@src/landing/components/PrizeInfo';
import { CountdownSection } from '@src/landing/components/CountdownSection';
import { Banner } from '@src/landing/components/Banner';
import { Footer } from '@src/components/Footer/Footer';
import { Actions } from '@src/landing/components/Actions';
import { Header } from '@src/components/Header/Header';

export const LandingPage = () => {
  return (
    <>
      <Header />
      <main className='max-w-5xl mx-auto p-4 bg-gray-900 text-white min-h-screen'>
        <PrizeInfo />
        <CountdownSection />
        <Banner />
        <Actions />
      </main>
      <Footer />
    </>
  );
};
