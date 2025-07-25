import React, { useState } from 'react';
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@src/components/ui/select';
import { Button } from '@src/components/ui/button';

type PurchaseFiltersProps = {
  onFilter: (filters: { status: string }) => void;
};

export const PurchaseFilters: React.FC<PurchaseFiltersProps> = ({
  onFilter,
}) => {
  const [status, setStatus] = useState<string>('');

  const handleApply = () => {
    onFilter({ status });
  };

  return (
    <div
      className='
        flex flex-col gap-3 p-4 bg-gray-900 rounded-2xl shadow mb-4
        sm:flex-row sm:items-center sm:gap-4
      '
    >
      <Select value={status} onValueChange={setStatus}>
        <SelectTrigger className='w-full sm:w-auto sm:min-w-[160px] sm:max-w-xs bg-gray-900 text-white border-gray-700'>
          <SelectValue placeholder='Seleccionar estado' />
        </SelectTrigger>
        <SelectContent className='bg-gray-900 text-white border-gray-700'>
          <SelectItem value='all'>Todos</SelectItem>
          <SelectItem value='pending'>Pendiente</SelectItem>
          <SelectItem value='verified'>Completado</SelectItem>
          <SelectItem value='cancelled'>Cancelado</SelectItem>
        </SelectContent>
      </Select>
      <Button
        onClick={handleApply}
        className='w-full sm:w-auto bg-gray-700 text-white hover:bg-gray-600'
      >
        Aplicar
      </Button>
    </div>
  );
};
