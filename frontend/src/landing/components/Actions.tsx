import { useAuth } from '@src/context/useAuth';
import { useModal } from '@src/context/useModal';
import { useState } from 'react';
import { PurchaseModal } from './PurchaseModal';
import { TicketsModal } from './TicketsModal';

export const Actions = () => {
  const { user } = useAuth();
  const { openLoginModal } = useModal();
  const [purchaseModalOpen, setPurchaseModalOpen] = useState(false);
  const [ticketModalOpen, setTicketModalOpen] = useState(false);

  const isLoggedIn = !!user;

  const handleCheckNumbersClick = () => {
    if (!isLoggedIn) {
      openLoginModal(false);
    } else {
      setTicketModalOpen(true);
    }
  };

  const handleBuyClick = () => {
    if (!isLoggedIn) {
      openLoginModal(false);
    } else {
      setPurchaseModalOpen(true);
    }
  };

  return (
    <>
      <div className='max-w-xs mx-auto mt-8 flex flex-col gap-4 px-4'>
        <button
          onClick={handleBuyClick}
          className='bg-[#272525] hover:bg-gray-800 text-white font-bold py-3 rounded transition ${
          isLoggedIn
            ? "bg-[#f30505] hover:bg-gray-800 cursor-pointer"
            : "bg-gray-600 cursor-not-allowed"
        }`'
        >
          🛒 COMPRAR AHORA
        </button>

        <PurchaseModal
          isOpen={purchaseModalOpen}
          onClose={() => setPurchaseModalOpen(false)}
        />
        <button
          onClick={handleCheckNumbersClick}
          className='bg-[#272525] hover:bg-gray-800 text-white font-bold py-3 rounded transition'
        >
          🔍 MIS BOLETOS COMPRADOS
        </button>
        {ticketModalOpen && user && (
          <TicketsModal
            userEmail={user.email}
            onClose={() => setTicketModalOpen(false)}
          />
        )}
      </div>
    </>
  );
};
