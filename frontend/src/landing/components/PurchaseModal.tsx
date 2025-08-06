import { useState } from 'react';
import { QuantitySelector } from './QuantitySelector';
import { PaymentMethods } from './PaymentMethods';
import { BuyForm } from './BuyForm';

interface PurchaseModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const steps = {
  QUANTITY: 'quantity',
  PAYMENT_METHOD: 'paymentMethod',
  BUY_FORM: 'buyForm',
} as const;

type PurchaseSteps = (typeof steps)[keyof typeof steps];

export const PurchaseModal = ({ isOpen, onClose }: PurchaseModalProps) => {
  const [step, setStep] = useState<PurchaseSteps>(steps.QUANTITY);
  const [quantity, setQuantity] = useState<number>(2);
  const [montoBs, setMontoBs] = useState<string>('150');
  const [montoUSD, setMontoUSD] = useState<string>('1');
  const [paymentMethod, setPaymentMethod] = useState<string | null>(null);
  const [selectedNumbers, setSelectedNumbers] = useState<string[]>([]);

  const handleQuantityNext = (
    qty: number,
    montoBs: string,
    montoUSD: string,
    numbers: string[],
  ) => {
    setQuantity(qty);
    setMontoBs(montoBs);
    setMontoUSD(montoUSD);
    setSelectedNumbers(numbers);
    setStep(steps.PAYMENT_METHOD);
  };

  const handleClose = () => {
    setStep(steps.QUANTITY);
    setPaymentMethod(null);
    onClose();
  };

  const handlePaymentBack = () => setStep(steps.QUANTITY);
  const handlePaymentNext = () => setStep(steps.BUY_FORM);
  const handleBuyFormBack = () => setStep(steps.PAYMENT_METHOD);
  if (!isOpen) return null;

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70'>
      <div className='bg-gray-900 p-6 rounded max-w-md w-full'>
        {step === steps.QUANTITY && (
          <QuantitySelector
            min={2}
            priceBS={150}
            priceUsd={1}
            onClose={handleClose}
            onNext={handleQuantityNext}
          />
        )}
        {step === steps.PAYMENT_METHOD && (
          <PaymentMethods
            selectedMethod={paymentMethod}
            setSelectedMethod={setPaymentMethod}
            onBack={handlePaymentBack}
            onClose={handleClose}
            onNext={handlePaymentNext}
          />
        )}
        {step === steps.BUY_FORM && (
          <BuyForm
            quantity={quantity}
            montoBs={montoBs}
            montoUSD={montoUSD}
            paymentMethod={paymentMethod ?? ''}
            selectedNumbers={selectedNumbers}
            onBack={handleBuyFormBack}
            onClose={handleClose}
          />
        )}
      </div>
    </div>
  );
};
