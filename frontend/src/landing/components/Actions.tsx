import { useAuth } from '@src/context/useAuth';

export const Actions = () => {
  const { user } = useAuth();
  const isLoggedIn = !!user;

  return (
    <div className='max-w-xs mx-auto mt-8 flex flex-col gap-4 px-4'>
      <button
        disabled={!isLoggedIn}
        onClick={() => console.log('click')}
        className='bg-[#272525] hover:bg-gray-800 text-white font-bold py-3 rounded transition ${
          isLoggedIn
            ? "bg-[#f30505] hover:bg-gray-800 cursor-pointer"
            : "bg-gray-600 cursor-not-allowed"
        }`'
      >
        🛒 COMPRAR AHORA
      </button>
      <button
        disabled={!isLoggedIn}
        className='bg-[#272525] hover:bg-gray-800 text-white font-bold py-3 rounded transition'
      >
        🔍 MIS NÚMEROS COMPRADOS
      </button>
      <button className='bg-[#272525] hover:bg-gray-800 text-white font-bold py-3 rounded transition'>
        ⚡ TICKETS PARTICIPANDO
      </button>
    </div>
  );
};
