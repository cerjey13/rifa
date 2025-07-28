import React, { useState, useEffect } from 'react';

type TicketApiError = {
  message: string;
  taken?: string[];
};

//TODO: change to a real backend integration
async function mockPurchaseTickets(numbers: string[]): Promise<void> {
  await new Promise((res) => setTimeout(res, 800));
  const taken = numbers.filter(
    (number) => number === '100' || number === '150',
  );
  if (taken.length) {
    throw { message: 'Algunos números no están disponibles', taken };
  }
}
const MIN_TICKETS = 2 as const;
const MAX_TICKETS = 500 as const;
const TICKET_MIN_VALUE = 0 as const;
const TICKET_MAX_VALUE = 9999 as const;
const TICKET_INPUT_MAXLEN = 4 as const;

interface QuantitySelectorProps {
  min: number;
  max?: number;
  priceBS: number;
  priceUsd: number;
  onClose: () => void;
  onNext: (
    quantity: number,
    montoBs: string,
    montoUSD: string,
    numbers: string[],
  ) => void;
}

function calculateMonto(
  quantity: number,
  priceBS: number,
  conversionRate: number,
): string {
  return (quantity * priceBS * conversionRate).toFixed(2);
}

export const QuantitySelector = ({
  min = MIN_TICKETS,
  max = MAX_TICKETS,
  priceBS,
  priceUsd,
  onClose,
  onNext,
}: QuantitySelectorProps) => {
  const [quantity, setQuantity] = useState<number>(min);
  const [montoUSD, setMontoUSD] = useState<string>(
    calculateMonto(min, priceUsd, 1),
  );
  const [montoBs, setMontoBs] = useState<string>(
    calculateMonto(min, priceBS, 1),
  );
  const [numbers, setNumbers] = useState<string[]>(Array(min).fill(''));
  const [inputErrors, setInputErrors] = useState<string[]>([]);
  const [backendErrors, setBackendErrors] = useState<string[]>([]);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [generalError, setGeneralError] = useState<string | null>(null);

  const handleNumberChange = (idx: number, value: string) => {
    const num = value.replace(/^0+(?!$)/, '');
    if (
      num === '' ||
      (/^\d{1,4}$/.test(num) &&
        Number(num) >= TICKET_MIN_VALUE &&
        Number(num) <= TICKET_MAX_VALUE)
    ) {
      setNumbers((prev) =>
        prev.map((prevNum, i) => (i === idx ? num : prevNum)),
      );
      setBackendErrors((prev) => prev.filter((errNum) => errNum !== num));
      setGeneralError(null);
    }
  };

  useEffect(() => {
    setNumbers((prev) => {
      if (quantity > prev.length) {
        return [...prev, ...Array(quantity - prev.length).fill('')];
      } else {
        return prev.slice(0, quantity);
      }
    });
    setInputErrors((prev) => {
      if (quantity > prev.length) {
        return [...prev, ...Array(quantity - prev.length).fill('')];
      } else {
        return prev.slice(0, quantity);
      }
    });
    setMontoUSD(calculateMonto(quantity, priceUsd, 1));
    setMontoBs(calculateMonto(quantity, priceBS, 1));
    setBackendErrors([]);
    setGeneralError(null);
  }, [quantity, priceBS, priceUsd]);

  useEffect(() => {
    const seen = new Set<string>();
    const errors = numbers.map((num) => {
      if (num === '') return '';
      if (isNaN(Number(num))) return 'Inválido';
      if (Number(num) < TICKET_MIN_VALUE || Number(num) > TICKET_MAX_VALUE)
        return `${TICKET_MIN_VALUE}-${TICKET_MAX_VALUE}`;
      if (seen.has(num)) return 'Repetido';
      seen.add(num);
      return '';
    });
    setInputErrors(errors);
  }, [numbers]);

  const increment = () => setQuantity((q) => Math.min(q + 1, max));
  const decrement = () => setQuantity((q) => Math.max(q - 1, min));
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let val = Number(e.target.value);
    if (isNaN(val)) val = min;
    val = Math.min(Math.max(val, min), max);
    setQuantity(val);
  };

  const allValid =
    numbers.length === quantity && inputErrors.every((err) => !err);

  const handleNext = async () => {
    setSubmitting(true);
    setBackendErrors([]);
    setGeneralError(null);

    const numbersToCheck = numbers.filter((num) => num !== '');

    try {
      if (numbersToCheck.length > 0) {
        await mockPurchaseTickets(numbersToCheck);
      }
      onNext(quantity, montoBs, montoUSD, numbers);
    } catch (error) {
      if (typeof error === 'object' && error && 'message' in error) {
        const err = error as TicketApiError;
        if (err.taken && Array.isArray(err.taken)) {
          setBackendErrors(err.taken);
        }
        setGeneralError(err.message || 'Ocurrió un error. Intenta de nuevo.');
      } else {
        setGeneralError('Ocurrió un error. Intenta de nuevo.');
      }
    }

    setSubmitting(false);
  };

  return (
    <div className='flex flex-col gap-6'>
      <h2 className='font-bold uppercase text-lg'>
        Elige la cantidad de números a comprar
      </h2>
      <hr className='border-gray-700' />
      <p className='font-semibold text-red-500'>
        Si no seleccionas los números, se asignarán aleatoriamente.
      </p>
      <hr className='border-gray-700' />
      <p className='text-yellow-400 text-sm'>
        La cantidad mínima de números es: ({min})
      </p>
      <div className='flex justify-between text-base font-semibold'>
        <span>Monto BS: {montoBs}</span>
        <span>Monto ($): {montoUSD}</span>
      </div>

      <div className='flex items-center justify-center gap-4'>
        <button
          onClick={decrement}
          disabled={quantity <= min || submitting}
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
          disabled={submitting}
        />

        <button
          onClick={increment}
          disabled={quantity >= max || submitting}
          className='bg-gray-800 rounded px-4 py-2 text-2xl disabled:opacity-50 disabled:cursor-not-allowed'
          aria-label='Aumentar cantidad'
        >
          +
        </button>
      </div>

      <div className='flex flex-wrap gap-3'>
        {numbers.map((num, idx) => (
          <div key={idx} className='flex flex-col items-center w-20'>
            <input
              type='text'
              inputMode='numeric'
              value={num}
              onChange={(e) => handleNumberChange(idx, e.target.value)}
              className={`w-full text-center bg-gray-700 rounded py-2 border
                ${
                  inputErrors[idx] || backendErrors.includes(num)
                    ? 'border-red-500'
                    : 'border-transparent'
                }
              `}
              placeholder={`#${idx + 1}`}
              aria-label={`Número ${idx + 1}`}
              maxLength={TICKET_INPUT_MAXLEN}
              disabled={submitting}
            />
            {inputErrors[idx] && (
              <span className='text-red-500 text-xs'>{inputErrors[idx]}</span>
            )}
            {!inputErrors[idx] && backendErrors.includes(num) && (
              <span className='text-red-500 text-xs'>No disponible</span>
            )}
          </div>
        ))}
      </div>

      {generalError && (
        <div className='text-red-500 text-center'>{generalError}</div>
      )}

      <div className='flex gap-4 justify-end'>
        <button
          onClick={onClose}
          className='bg-gray-700 py-2 px-6 rounded hover:bg-gray-600 transition'
          disabled={submitting}
        >
          Cancelar
        </button>
        <button
          onClick={handleNext}
          className='bg-yellow-500 py-2 px-6 rounded hover:bg-yellow-600 text-black font-bold transition'
          disabled={!allValid || submitting}
        >
          {submitting ? 'Verificando...' : 'Siguiente'}
        </button>
      </div>
    </div>
  );
};
