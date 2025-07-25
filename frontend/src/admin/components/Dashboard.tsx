import { useState } from 'react';
import { ResumenCompras } from './AllBrief';
import { TicketsLeaderboard } from './MostTickets';

export const Dashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'resumen' | 'mascomprados'>(
    'resumen',
  );

  return (
    <div className='p-2 sm:p-4 md:p-6'>
      <div className='flex gap-2 mb-6'>
        <button
          className={`px-4 py-2 rounded-t-lg font-semibold
            ${
              activeTab === 'resumen'
                ? 'bg-gray-800 text-white border-b-2 border-blue-500'
                : 'bg-gray-700 text-gray-300'
            }`}
          onClick={() => setActiveTab('resumen')}
        >
          Resumen de compras
        </button>
        <button
          className={`px-4 py-2 rounded-t-lg font-semibold
            ${
              activeTab === 'mascomprados'
                ? 'bg-gray-800 text-white border-b-2 border-blue-500'
                : 'bg-gray-700 text-gray-300'
            }`}
          onClick={() => setActiveTab('mascomprados')}
        >
          Más comprados
        </button>
      </div>

      {activeTab === 'resumen' && <ResumenCompras />}

      {activeTab === 'mascomprados' && <TicketsLeaderboard />}
    </div>
  );
};
