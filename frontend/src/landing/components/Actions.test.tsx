import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi, describe, it, beforeEach, expect } from 'vitest';
import Actions, { type ActionsProps } from './Actions';

// Mock config fallbacks
vi.mock('@src/config/config', () => ({
  MONTO_BS: 123,
  MONTO_USD: 9,
}));

// Mock auth/modal hooks
const openLoginModal = vi.fn<(b: boolean) => void>();
vi.mock('@src/context/useModal', () => ({
  useModal: () => ({ openLoginModal }),
}));

const useAuthMock = vi.fn<() => { user: { email: string } | null }>();
vi.mock('@src/context/useAuth', () => ({
  useAuth: () => useAuthMock(),
}));

// Mock child modals (light)
vi.mock('./PurchaseModal', () => ({
  PurchaseModal: ({
    isOpen,
    bs,
    usd,
  }: {
    isOpen: boolean;
    bs: number;
    usd: number;
  }) =>
    isOpen ? <div data-testid='purchase' data-bs={bs} data-usd={usd} /> : null,
}));

vi.mock('./TicketsModal', () => ({
  TicketsModal: ({ userEmail }: { userEmail: string }) => (
    <div data-testid='tickets' data-email={userEmail} />
  ),
}));

function renderActions(props: ActionsProps) {
  return render(<Actions {...props} />);
}

describe('<Actions />', () => {
  beforeEach(() => {
    openLoginModal.mockClear();
    useAuthMock.mockReset();
  });

  it('logged out: clicking buy or check opens login modal', async () => {
    useAuthMock.mockReturnValue({ user: null });
    renderActions({ prices: undefined });

    await userEvent.click(
      screen.getByRole('button', { name: /COMPRAR AHORA/i }),
    );
    await userEvent.click(
      screen.getByRole('button', { name: /MIS BOLETOS COMPRADOS/i }),
    );
    expect(openLoginModal).toHaveBeenCalledTimes(2);
    expect(openLoginModal).toHaveBeenNthCalledWith(1, false);
    expect(openLoginModal).toHaveBeenNthCalledWith(2, false);

    // no modals rendered
    expect(screen.queryByTestId('purchase')).not.toBeInTheDocument();
    expect(screen.queryByTestId('tickets')).not.toBeInTheDocument();
  });

  it('logged in: clicking buy opens PurchaseModal with fallback prices when props undefined', async () => {
    useAuthMock.mockReturnValue({ user: { email: 'a@b.com' } });
    renderActions({ prices: undefined });

    await userEvent.click(
      screen.getByRole('button', { name: /COMPRAR AHORA/i }),
    );
    const modal = await screen.findByTestId('purchase');
    expect(modal.dataset.bs).toBe('123');
    expect(modal.dataset.usd).toBe('9');
  });

  it('logged in: clicking check opens TicketsModal with user email', async () => {
    useAuthMock.mockReturnValue({ user: { email: 'me@site.com' } });
    renderActions({ prices: { montoBs: 350, montoUsd: 1 } });

    await userEvent.click(
      screen.getByRole('button', { name: /MIS BOLETOS COMPRADOS/i }),
    );
    const tm = await screen.findByTestId('tickets');
    expect(tm.dataset.email).toBe('me@site.com');
  });

  it('uses provided prices when available', async () => {
    useAuthMock.mockReturnValue({ user: { email: 'x@y.z' } });
    renderActions({ prices: { montoBs: 350, montoUsd: 2 } });

    await userEvent.click(
      screen.getByRole('button', { name: /COMPRAR AHORA/i }),
    );
    const modal = await screen.findByTestId('purchase');
    expect(modal.dataset.bs).toBe('350');
    expect(modal.dataset.usd).toBe('2');
  });
});
