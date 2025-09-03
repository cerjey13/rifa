import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PurchaseModal } from './PurchaseModal';

vi.mock('@src/api/tickets', () => ({
  fetchAvailableTickets: vi.fn<() => Promise<string[]>>(),
}));

describe('<PurchaseModal /> integration (real QS + Payment)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders null when closed', () => {
    const { container } = render(
      <PurchaseModal bs={100} usd={10} isOpen={false} onClose={() => {}} />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('flows Quantity -> PaymentMethods -> BuyForm with default qty (2)', async () => {
    render(<PurchaseModal bs={100} usd={10} isOpen onClose={() => {}} />);
    // Step 1
    const next1 = screen.getByRole('button', { name: 'Siguiente' });
    expect(next1).toBeEnabled();
    await userEvent.click(next1);

    // Step 2
    const pagoMovil = await screen.findByRole('button', { name: /PAGOMOVIL/i });
    await userEvent.click(pagoMovil);
    const next2 = screen.getByRole('button', { name: 'Siguiente' });
    expect(next2).toBeEnabled();
    await userEvent.click(next2);

    // Step 3
    const quantity = screen.getByLabelText('cantidad') as HTMLSpanElement;
    const bs = screen.getByLabelText('bs') as HTMLSpanElement;
    expect(quantity.textContent).toContain('2');
    expect(bs.textContent).toContain('200');
  });

  it('allows changing quantity (to 3) before proceeding and passes updated amounts', async () => {
    render(<PurchaseModal bs={100} usd={10} isOpen onClose={() => {}} />);

    const inc = screen.getByLabelText('Aumentar cantidad');
    await userEvent.click(inc);

    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    await userEvent.click(
      await screen.findByRole('button', { name: /ZELLE/i }),
    );
    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    const quantity = screen.getByLabelText('cantidad') as HTMLSpanElement;
    const usd = screen.getByLabelText('usd') as HTMLSpanElement;
    expect(quantity.textContent).toContain('3');
    expect(usd.textContent).toContain('30');
  });

  it('PaymentMethods: Next disabled until selected; toggling off disables again', async () => {
    render(<PurchaseModal bs={100} usd={10} isOpen onClose={() => {}} />);

    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    const next = await screen.findByRole('button', { name: 'Siguiente' });
    expect(next).toBeDisabled();

    const pagoMovil = screen.getByRole('button', { name: /PAGOMOVIL/i });
    await userEvent.click(pagoMovil);
    expect(next).toBeEnabled();

    await userEvent.click(pagoMovil);
    expect(next).toBeDisabled();
  });

  it('PaymentMethods: "Atrás" returns to Quantity step', async () => {
    render(<PurchaseModal bs={100} usd={10} isOpen onClose={() => {}} />);

    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    await userEvent.click(screen.getByRole('button', { name: 'Atrás' }));

    expect(
      screen.getByRole('button', { name: 'Siguiente' }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('Cantidad de números')).toBeInTheDocument();
  });

  it('PaymentMethods: "Cerrar" calls onClose and resets internal step to Quantity', async () => {
    const onClose = vi.fn();
    render(<PurchaseModal bs={100} usd={10} isOpen onClose={onClose} />);

    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    const closeBtn = await screen.findByRole('button', { name: 'Cerrar' });
    await userEvent.click(closeBtn);

    expect(onClose).toHaveBeenCalledTimes(1);
    expect(
      screen.getByRole('button', { name: 'Siguiente' }),
    ).toBeInTheDocument();
  });

  it('BuyForm: "bf-back" returns to Payment; "bf-close" calls onClose and resets to Quantity', async () => {
    const onClose = vi.fn();
    render(<PurchaseModal bs={100} usd={10} isOpen onClose={onClose} />);

    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));
    await userEvent.click(
      await screen.findByRole('button', { name: /PAGOMOVIL/i }),
    );
    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    await userEvent.click(screen.getByRole('button', { name: 'Atrás' }));
    expect(
      await screen.findByRole('button', { name: /PAGOMOVIL/i }),
    ).toBeInTheDocument();
    await userEvent.click(
      await screen.findByRole('button', { name: /PAGOMOVIL/i }),
    );
    await userEvent.click(screen.getByRole('button', { name: 'Siguiente' }));

    await userEvent.click(screen.getByRole('button', { name: 'Cerrar' }));
    expect(onClose).toHaveBeenCalledTimes(1);
    expect(
      screen.getByRole('button', { name: 'Siguiente' }),
    ).toBeInTheDocument();
  });
});
