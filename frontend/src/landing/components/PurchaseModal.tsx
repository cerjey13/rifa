import { useState } from 'react';
import { QuantitySelector } from './QuantitySelector';
import { PaymentMethods } from './PaymentMethods';
import { BuyForm } from './BuyForm';

interface PurchaseModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const PurchaseModal = ({ isOpen, onClose }: PurchaseModalProps) => {
  const [step, setStep] = useState<'quantity' | 'paymentMethod' | 'buyForm'>(
    'quantity',
  );
  const [quantity, setQuantity] = useState<number>(2);
  const [montoBs, setMontoBs] = useState<string>('40');
  const [montoUSD, setMontoUSD] = useState<string>('');
  const [paymentMethod, setPaymentMethod] = useState<string | null>(null);
  const handleQuantityNext = (
    qty: number,
    montoBs: string,
    montoUSD: string,
  ) => {
    setQuantity(qty);
    setMontoBs(montoBs);
    setMontoUSD(montoUSD);
    setStep('paymentMethod');
  };

  const handleClose = () => {
    setStep('quantity');
    setPaymentMethod(null);
    onClose();
  };

  const handlePaymentBack = () => setStep('quantity');
  const handlePaymentNext = () => setStep('buyForm');
  const handleBuyFormBack = () => setStep('paymentMethod');
  if (!isOpen) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70'>
      <div className='bg-gray-900 p-6 rounded max-w-md w-full'>
        {step === 'quantity' && (
          <QuantitySelector
            min={2}
            priceBS={40}
            conversionRate={0.009}
            onClose={handleClose}
            onNext={handleQuantityNext}
          />
        )}
        {step === 'paymentMethod' && (
          <PaymentMethods
            selectedMethod={paymentMethod}
            setSelectedMethod={setPaymentMethod}
            onBack={handlePaymentBack}
            onClose={handleClose}
            onNext={handlePaymentNext}
          />
        )}
        {step === 'buyForm' && (
          <BuyForm
            quantity={quantity}
            montoBs={montoBs}
            montoUSD={montoUSD}
            paymentMethod={paymentMethod ?? ''}
            onBack={handleBuyFormBack}
            onClose={handleClose}
          />
        )}
      </div>
    </div>
  );
};
