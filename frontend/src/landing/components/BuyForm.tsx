import { useState } from 'react';
import type { ChangeEvent, FormEvent } from 'react';
import { submitPurchase } from '@src/api/purchase';
import { CopyableText } from '@src/components/Clipboard/Copy';
import { toast } from 'sonner';

interface BuyFormProps {
  quantity: number;
  montoBs: number;
  montoUSD: number;
  paymentMethod: string;
  selectedNumbers: string[];
  onBack: () => void;
  onClose: () => void;
}

const paymentMethods = {
  pagoMovil: 'pago movil',
  zelle: 'zelle',
} as const;

export const BuyForm = ({
  quantity,
  montoBs,
  montoUSD,
  paymentMethod,
  selectedNumbers,
  onBack,
  onClose,
}: BuyFormProps) => {
  const [transactionDigits, setTransactionDigits] = useState('');
  const [paymentScreenshot, setPaymentScreenshot] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file =
      e.target.files && e.target.files.length > 0 ? e.target.files[0] : null;
    if (file) {
      const maxSize = 3 * 1024 * 1024; // 3MB in bytes
      if (file.size > maxSize) {
        toast.error('El archivo no debe superar los 3 MB.', { duration: 5000 });
        e.target.value = ''; // reset input
        return;
      }
      setPaymentScreenshot(file);
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    if (transactionDigits.length !== 6) {
      alert('Por favor ingresa los últimos 6 dígitos de la transacción');
      return;
    }
    if (!paymentScreenshot) {
      alert('Por favor sube una captura del pago realizado');
      return;
    }

    setLoading(true);

    try {
      await submitPurchase({
        quantity,
        montoBs: montoBs.toFixed(2),
        montoUSD: montoUSD.toFixed(2),
        paymentMethod,
        transactionDigits,
        selectedNumbers,
        paymentScreenshot,
      });

      setSuccess(true);
    } catch {
      setError('Error al enviar los datos, intenta nuevamente');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div
        className={`fixed inset-0 flex items-center justify-center bg-black bg-opacity-70 z-50 transition-opacity duration-500 ${
          success
            ? 'opacity-100 pointer-events-auto'
            : 'opacity-0 pointer-events-none'
        }`}
        aria-hidden={!success}
      >
        <div
          className="bg-gray-800 p-8 rounded-md text-white max-w-sm w-full text-center transition-transform duration-500 ease-out transform ${
          success ? 'translate-y-0' : '-translate-y-10'
        }"
        >
          <h2 className='text-xl font-bold mb-4'>¡Éxito!</h2>
          <p>Tu pago fue registrado correctamente.</p>
          <button
            onClick={onClose}
            className='mt-6 bg-yellow-500 hover:bg-yellow-600 text-black font-bold py-2 px-4 rounded'
          >
            Cerrar
          </button>
        </div>
      </div>

      {!success && (
        <form
          onSubmit={handleSubmit}
          className='flex flex-col gap-6 text-white  max-h-[90vh] overflow-y-auto py-4'
        >
          <div className='flex items-center justify-between'>
            <h2 className='font-bold uppercase text-lg'>Pagar</h2>

            <button
              onClick={onClose}
              aria-label='Cerrar'
              className='p-2 text-white hover:text-yellow-400 transition'
            >
              <svg
                xmlns='http://www.w3.org/2000/svg'
                className='h-6 w-6'
                fill='none'
                viewBox='0 0 24 24'
                stroke='currentColor'
                strokeWidth={2}
              >
                <path
                  strokeLinecap='round'
                  strokeLinejoin='round'
                  d='M6 18L18 6M6 6l12 12'
                />
              </svg>
            </button>
          </div>
          <hr className='border-gray-700' />

          <p className='text-sm mb-4'>Envia tu registro de pago y participa</p>

          <div className='space-y-1'>
            <p>
              <span className='font-semibold'>Metodo:</span>{' '}
              {paymentMethod.toUpperCase()}
            </p>
            {paymentMethod === paymentMethods.pagoMovil && (
              <>
                <p>
                  <span className='font-semibold'>Cuenta:</span>{' '}
                  <CopyableText text='BANESCO 0134' copyText='BANESCO' />
                </p>
                <p>
                  <span className='font-semibold'>Cédula:</span>{' '}
                  <CopyableText text='30606459' />
                </p>
                <p>
                  <span className='font-semibold'>Teléfono:</span>{' '}
                  <CopyableText text='04141551801' />
                </p>
              </>
            )}

            {paymentMethod === paymentMethods.zelle && (
              <>
                <p>
                  <span className='font-semibold'>Teléfono:</span>{' '}
                  <CopyableText text='3802389306' />
                </p>
                <p>
                  <span className='font-semibold'>Nombre:</span>{' '}
                  <CopyableText text='Vicente Méndez' />
                </p>
              </>
            )}
          </div>

          <hr className='border-gray-700' />

          <div className='space-y-1'>
            <p>
              <span className='font-semibold'>Cantidad de tickets:</span>{' '}
              {quantity}
            </p>
            {paymentMethod === paymentMethods.pagoMovil && (
              <p>
                <span className='font-semibold'>Monto BS:</span> {montoBs}
              </p>
            )}
            {paymentMethod === paymentMethods.zelle && (
              <p>
                <span className='font-semibold'>Monto ($):</span> {montoUSD}
              </p>
            )}
            {selectedNumbers.filter(Boolean).length > 0 && (
              <p>
                <span className='font-semibold'>Números seleccionados:</span>{' '}
                {selectedNumbers.filter(Boolean).join(', ')}
              </p>
            )}
          </div>

          <hr className='border-gray-700' />

          <label className='flex flex-col gap-1'>
            Últimos 6 Dígitos (Transacción)
            <input
              type='text'
              maxLength={6}
              pattern='\d{6}'
              required
              value={transactionDigits}
              onChange={(e) =>
                setTransactionDigits(e.target.value.replace(/\D/g, ''))
              }
              className='bg-gray-700 rounded mx-0.5 p-2 text-white focus:outline-none focus:ring-2 focus:ring-yellow-500'
              placeholder='123456'
              disabled={loading}
            />
          </label>

          <div>
            <p className='font-semibold mb-1'>Captura del pago realizado</p>
            <input
              name='file'
              type='file'
              accept='image/*'
              onChange={handleFileChange}
              className='w-full cursor-pointer text-sm text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-yellow-500 file:text-black hover:file:bg-yellow-600'
              required
              disabled={loading}
            />
            {paymentScreenshot && (
              <p className='mt-2 text-sm text-gray-300'>
                Archivo seleccionado: {paymentScreenshot.name}
              </p>
            )}
          </div>

          {error && <p className='text-red-500 text-sm'>{error}</p>}

          <div className='flex gap-4 justify-end'>
            <button
              type='button'
              onClick={onBack}
              className='bg-gray-700 py-2 px-6 rounded hover:bg-gray-600 transition'
            >
              Atrás
            </button>
            <button
              type='submit'
              disabled={loading}
              aria-busy={loading}
              className='bg-yellow-500 py-2 px-6 rounded hover:bg-yellow-600 text-black font-bold transition'
            >
              {loading && (
                <svg
                  className='animate-spin h-5 w-5 text-black'
                  xmlns='http://www.w3.org/2000/svg'
                  fill='none'
                  viewBox='0 0 24 24'
                >
                  <circle
                    className='opacity-25'
                    cx='12'
                    cy='12'
                    r='10'
                    stroke='currentColor'
                    strokeWidth='4'
                  />
                  <path
                    className='opacity-75'
                    fill='currentColor'
                    d='M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z'
                  />
                </svg>
              )}
              {!loading && 'Finalizar'}
            </button>
          </div>
        </form>
      )}
    </>
  );
};
