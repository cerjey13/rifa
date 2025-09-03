import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { PaymentMethods } from './PaymentMethods';
import { useState } from 'react';

function Wrapped({
  onBack = () => {},
  onClose = () => {},
  onNext = () => {},
}: {
  onBack?: () => void;
  onClose?: () => void;
  onNext?: () => void;
}) {
  const [method, setMethod] = useState<string | null>(null);
  return (
    <PaymentMethods
      selectedMethod={method}
      setSelectedMethod={setMethod}
      onBack={onBack}
      onClose={onClose}
      onNext={onNext}
    />
  );
}

describe('<PaymentMethods />', () => {
  it('Next disabled until a method is selected, enabled after select and disabled again if toggled off', async () => {
    render(<Wrapped />);
    const next = screen.getByRole('button', { name: /Siguiente/i });
    expect(next).toBeDisabled();

    const pagoMovil = screen.getByRole('button', { name: /PAGOMOVIL/i });
    await userEvent.click(pagoMovil);
    expect(next).toBeEnabled();

    // toggle off
    await userEvent.click(pagoMovil);
    expect(next).toBeDisabled();
  });

  it('Close clears selection and calls onClose', async () => {
    const onClose = vi.fn();
    render(<Wrapped onClose={onClose} />);

    // select first to ensure it clears
    await userEvent.click(screen.getByRole('button', { name: /PAGOMOVIL/i }));

    const closeBtn = screen.getByRole('button', { name: /Cerrar/i });
    await userEvent.click(closeBtn);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('Back and Next call handlers', async () => {
    const onBack = vi.fn();
    const onNext = vi.fn();
    render(<Wrapped onBack={onBack} onNext={onNext} />);

    // Back
    await userEvent.click(screen.getByRole('button', { name: /Atrás/i }));
    expect(onBack).toHaveBeenCalledTimes(1);

    // Select and Next
    await userEvent.click(screen.getByRole('button', { name: /ZELLE/i }));
    const next = screen.getByRole('button', { name: /Siguiente/i });
    await userEvent.click(next);
    expect(onNext).toHaveBeenCalledTimes(1);
  });
});
