import { useState } from 'react';
import { ResumenCompras } from '@src/admin/components/AllBrief';
import { TicketsLeaderboard } from '@src/admin/components/MostTickets';
import { Searcher } from '@src/admin/components/Searcher';

export const Dashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<
    'resumen' | 'mascomprados' | 'buscador'
  >('resumen');

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
        <button
          className={`px-4 py-2 rounded-t-lg font-semibold
            ${
              activeTab === 'buscador'
                ? 'bg-gray-800 text-white border-b-2 border-blue-500'
                : 'bg-gray-700 text-gray-300'
            }`}
          onClick={() => setActiveTab('buscador')}
        >
          Buscador
        </button>
      </div>

      {activeTab === 'resumen' && <ResumenCompras />}

      {activeTab === 'mascomprados' && <TicketsLeaderboard />}

      {activeTab === 'buscador' && <Searcher />}
    </div>
  );
};
