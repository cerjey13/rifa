import React, { useState } from 'react';

interface QuantitySelectorProps {
  min: number;
  max?: number;
  priceBS: number;
  priceUsd: number;
  onClose: () => void;
  onNext: (quantity: number, montoBs: string, montoUSD: string) => void;
}

function calculateMonto(
  quantity: number,
  priceBS: number,
  conversionRate: number,
): string {
  return (quantity * priceBS * conversionRate).toFixed(2);
}

export const QuantitySelector = ({
  min = 2,
  max = 500,
  priceBS,
  priceUsd,
  onClose,
  onNext,
}: QuantitySelectorProps) => {
  const [quantity, setQuantity] = useState<number>(min);
  const [montoUSD, setMontoUSD] = useState<string>(
    calculateMonto(quantity, priceUsd, 1),
  );
  const [montoBs, setMontoBs] = useState<string>(
    calculateMonto(quantity, priceBS, 1),
  );

  const increment = () => {
    setQuantity((q) => Math.min(q + 1, max));
    setMontoUSD(calculateMonto(quantity + 1, priceUsd, 1));
    setMontoBs(calculateMonto(quantity + 1, priceBS, 1));
  };
  const decrement = () => {
    setQuantity((q) => Math.max(q - 1, min));
    setMontoUSD(calculateMonto(quantity - 1, priceUsd, 1));
    setMontoBs(calculateMonto(quantity - 1, priceBS, 1));
  };
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let val = Number(e.target.value);
    if (isNaN(val)) val = min;
    val = Math.min(Math.max(val, min), max);
    setQuantity(val);
  };

  return (
    <div className='flex flex-col gap-6'>
      <h2 className='font-bold uppercase text-lg'>
        Elige la cantidad de números a comprar
      </h2>
      <hr className='border-gray-700' />
      <p className='text-yellow-400 text-sm'>
        La cantidad mínima de números es: ({min})
      </p>
      <div className='flex justify-between text-base font-semibold'>
        <span>Monto BS: {calculateMonto(quantity, priceBS, 1)}</span>
        <span>Monto ($): {calculateMonto(quantity, priceUsd, 1)}</span>
      </div>

      <div className='flex items-center justify-center gap-4'>
        <button
          onClick={decrement}
          disabled={quantity <= min}
          className='bg-gray-800 rounded px-4 py-2 text-2xl disabled:opacity-50 disabled:cursor-not-allowed'
          aria-label='Disminuir cantidad'
        >
          -
        </button>

        <input
          type='number'
          min={min}
          max={max}
          value={quantity}
          onChange={handleChange}
          className='w-20 text-center bg-gray-700 rounded py-2 focus:outline-none focus:ring-2 focus:ring-yellow-500'
          aria-label='Cantidad de números'
        />

        <button
          onClick={increment}
          disabled={quantity >= max}
          className='bg-gray-800 rounded px-4 py-2 text-2xl disabled:opacity-50 disabled:cursor-not-allowed'
          aria-label='Aumentar cantidad'
        >
          +
        </button>
      </div>

      <div className='flex gap-4 justify-end'>
        <button
          onClick={onClose}
          className='bg-gray-700 py-2 px-6 rounded hover:bg-gray-600 transition'
        >
          Cancelar
        </button>
        <button
          onClick={() => onNext(quantity, montoBs, montoUSD)}
          className='bg-yellow-500 py-2 px-6 rounded hover:bg-yellow-600 text-black font-bold transition'
        >
          Siguiente
        </button>
      </div>
    </div>
  );
};
