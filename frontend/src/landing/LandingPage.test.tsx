import { renderWithQuery } from '@src/test-utils';
import { screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';

// Mock leaf UI to make prop assertions easy (works with React.lazy)
vi.mock('@src/landing/components/CountdownSection', () => ({
  default: ({ percentage, loading, error }: CountdownSectionProps) => (
    <div
      data-testid='countdown'
      data-percentage={percentage}
      data-loading={loading}
      data-error={error}
    />
  ),
}));

vi.mock('@src/landing/components/Actions', () => ({
  default: ({ prices }: ActionsProps) => (
    <div data-testid='actions' data-prices={JSON.stringify(prices)} />
  ),
}));

// Keep header/footer/prize simple
vi.mock('@src/components/Header/Header', () => ({
  Header: () => <header>header</header>,
}));
vi.mock('@src/components/Footer/Footer', () => ({
  Footer: () => <footer>footer</footer>,
}));
vi.mock('@src/landing/components/PrizeInfo', () => ({
  PrizeInfo: () => <section>prize</section>,
}));

// Mock data layer
vi.mock('@src/api/tickets', () => ({
  fetchTicketsPercentage: vi.fn<() => Promise<number>>(),
}));
vi.mock('@src/api/prices', () => ({
  fetchPrices: vi.fn<() => Promise<Prices>>(),
}));

import { fetchTicketsPercentage } from '@src/api/tickets';
import { fetchPrices } from '@src/api/prices';
import { LandingPage } from './LandingPage';
import type { ActionsProps } from './components/Actions';
import type { CountdownSectionProps } from './components/CountdownSection';

describe('<LandingPage />', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders success path: passes percentage and prices to children', async () => {
    const ticketsMock = vi.mocked(fetchTicketsPercentage);
    const pricesMock = vi.mocked(fetchPrices);
    ticketsMock.mockResolvedValue(42);
    pricesMock.mockResolvedValue({ montoUsd: 1, montoBs: 350 });

    renderWithQuery(<LandingPage />);

    const countdown = await screen.findByTestId('countdown');
    const actions = await screen.findByTestId('actions');

    // CountdownSection props
    expect(countdown.dataset.percentage).toBe('42');
    expect(countdown.dataset.loading).toBe('false');
    expect(countdown.dataset.error).toBe('false');

    // Actions props
    await waitFor(() =>
      expect(actions.dataset.prices).toBe(
        JSON.stringify({ montoUsd: 1, montoBs: 350 }),
      ),
    );
  });

  it('handles ticket error: percentage=0, error=true; prices still shown if ok', async () => {
    const ticketsMock = vi.mocked(fetchTicketsPercentage);
    const pricesMock = vi.mocked(fetchPrices);

    ticketsMock.mockRejectedValue(new Error('boom'));
    pricesMock.mockResolvedValue({ montoUsd: 11, montoBs: 150 });

    renderWithQuery(<LandingPage />);

    const countdown = await screen.findByTestId('countdown');
    const actions = await screen.findByTestId('actions');

    expect(countdown.dataset.percentage).toBe('0');
    expect(countdown.dataset.loading).toBe('false');
    expect(countdown.dataset.error).toBe('true');

    expect(actions.dataset.prices).toBe(
      JSON.stringify({ montoUsd: 11, montoBs: 150 }),
    );
  });

  it('still renders when tickets and prices fail', async () => {
    const ticketsMock = vi.mocked(fetchTicketsPercentage);
    const pricesMock = vi.mocked(fetchPrices);

    ticketsMock.mockRejectedValue(new Error('tickets failed'));
    pricesMock.mockRejectedValue(new Error('prices failed'));

    renderWithQuery(<LandingPage />);

    const countdown = await screen.findByTestId('countdown');
    const actions = await screen.findByTestId('actions');

    // CountdownSection gets safe fallbacks
    expect(countdown.dataset.percentage).toBe('0');
    expect(countdown.dataset.loading).toBe('false');
    expect(countdown.dataset.error).toBe('true');

    // Actions still renders
    expect(actions).toBeInTheDocument();
    expect(actions.dataset.prices).toBeUndefined();
  });
});
