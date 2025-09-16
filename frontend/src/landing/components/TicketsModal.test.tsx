import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import { renderWithQuery } from '@src/test-utils';
import { TicketsModal } from './TicketsModal';

vi.mock('@src/api/tickets', () => ({
  fetchUserTickets: vi.fn<() => Promise<string[]>>(),
}));
import { fetchUserTickets } from '@src/api/tickets';

describe('<TicketsModal />', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns null (renders nothing) when userEmail is empty and does not fetch', () => {
    const { container } = renderWithQuery(
      <TicketsModal userEmail='' onClose={() => {}} />,
    );
    expect(container.firstChild).toBeNull();
    expect(vi.mocked(fetchUserTickets)).not.toHaveBeenCalled();
  });

  it('shows loading state while fetching', () => {
    // never-resolving promise to keep loading
    vi.mocked(fetchUserTickets).mockImplementation(() => new Promise(() => {}));

    renderWithQuery(<TicketsModal userEmail='a@b.com' onClose={() => {}} />);

    expect(screen.getByText(/Cargando\.\.\./i)).toBeInTheDocument();
  });

  it('shows error state on fetch failure', async () => {
    vi.mocked(fetchUserTickets).mockRejectedValue(new Error('boom'));

    renderWithQuery(<TicketsModal userEmail='a@b.com' onClose={() => {}} />);

    expect(
      await screen.findByText(/Error al cargar tickets/i),
    ).toBeInTheDocument();
  });

  it('shows empty state when there are no tickets', async () => {
    vi.mocked(fetchUserTickets).mockResolvedValue([]);

    renderWithQuery(<TicketsModal userEmail='a@b.com' onClose={() => {}} />);

    expect(
      await screen.findByText(/No has comprado ningún número\./i),
    ).toBeInTheDocument();
  });

  it('lists tickets and displays the count', async () => {
    const tickets = ['0001', '0100', '1024', '2048'];
    vi.mocked(fetchUserTickets).mockResolvedValue(tickets);

    renderWithQuery(<TicketsModal userEmail='a@b.com' onClose={() => {}} />);

    // heading appears
    expect(
      await screen.findByRole('heading', { name: /Mis números comprados/i }),
    ).toBeInTheDocument();

    // count text
    expect(
      screen.getByText(String(tickets.length), { selector: 'strong' }),
    ).toBeInTheDocument();

    // each ticket appears
    for (const t of tickets) {
      expect(screen.getByText(t)).toBeInTheDocument();
    }
  });

  it('calls onClose when the close button is clicked', async () => {
    vi.mocked(fetchUserTickets).mockResolvedValue(['0001']);

    const onClose = vi.fn();
    renderWithQuery(<TicketsModal userEmail='a@b.com' onClose={onClose} />);

    const btn = await screen.findByRole('button', { name: /Cerrar/i });
    await userEvent.click(btn);

    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
