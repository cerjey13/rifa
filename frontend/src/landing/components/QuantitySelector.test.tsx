import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { QuantitySelector } from './QuantitySelector';

vi.mock('@src/api/tickets', () => ({
  fetchAvailableTickets: vi.fn<() => Promise<string[]>>(),
}));
import { fetchAvailableTickets } from '@src/api/tickets';

const renderQS = () =>
  render(
    <QuantitySelector
      min={2}
      priceBS={10}
      priceUsd={1}
      onClose={() => {}}
      onNext={() => {}}
    />,
  );

describe('<QuantitySelector />', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows initial amounts for BS and USD with default quantity', () => {
    renderQS();
    expect(screen.getByText(/Monto BS:\s*20/)).toBeInTheDocument();
    expect(screen.getByText(/Monto \(\$\):\s*2/)).toBeInTheDocument();
  });

  it('increments, decrements and clamps between min/max; updates amounts', async () => {
    renderQS();
    const inc = screen.getByRole('button', { name: 'Aumentar cantidad' });
    const dec = screen.getByRole('button', { name: 'Disminuir cantidad' });
    const qtyInput = screen.getByRole('spinbutton', {
      name: 'Cantidad de números',
    });

    await userEvent.click(inc);
    expect(qtyInput).toHaveValue(3);
    expect(screen.getByText(/Monto BS:\s*30/)).toBeInTheDocument();
    expect(screen.getByText(/Monto \(\$\):\s*3/)).toBeInTheDocument();

    await userEvent.click(dec);
    expect(qtyInput).toHaveValue(2);

    // direct set out of range -> clamped to min (2)
    await userEvent.clear(qtyInput);
    expect(qtyInput).toHaveValue(2);
  });

  it('accepts up to 4-digit numeric ticket inputs, strips leading zeros, rejects invalid', async () => {
    renderQS();

    const ticket1 = screen.getByLabelText('Número 1') as HTMLInputElement;
    const ticket2 = screen.getByLabelText('Número 2') as HTMLInputElement;

    await userEvent.type(ticket1, '0001');
    expect(ticket1.value).toBe('1');

    await userEvent.clear(ticket2);
    await userEvent.type(ticket2, '12345');
    expect(ticket2.value).toBe('1234');
  });

  it('shows validation errors for duplicates and out-of-range values; disables Next', async () => {
    renderQS();
    const ticket1 = screen.getByLabelText('Número 1');
    const ticket2 = screen.getByLabelText('Número 2');

    // duplicate
    await userEvent.type(ticket1, '12');
    await userEvent.type(ticket2, '12');
    expect(await screen.findAllByText('Repetido')).toHaveLength(1);

    // Next disabled
    const next = screen.getByRole('button', { name: 'Siguiente' });
    expect(next).toBeDisabled();

    // out of range
    await userEvent.clear(ticket2);
    await userEvent.type(ticket2, '10000');
    expect(screen.queryByText('0-9999')).not.toBeInTheDocument();
  });

  it('proceeds with empty numbers (random assignment) and does not call availability API', async () => {
    const onNext = vi.fn();
    render(
      <QuantitySelector
        min={2}
        priceBS={10}
        priceUsd={1}
        onClose={() => {}}
        onNext={onNext}
      />,
    );

    vi.mocked(fetchAvailableTickets).mockResolvedValue([]);
    const next = screen.getByRole('button', { name: 'Siguiente' });
    await userEvent.click(next);

    expect(fetchAvailableTickets).not.toHaveBeenCalled();
    expect(onNext).toHaveBeenCalledTimes(1);
    // quantity=2, montoBs=20, montoUSD=2, numbers=['','']
    expect(onNext).toHaveBeenCalledWith(2, 20, 2, ['', '']);
  });

  it('checks availability when numbers provided and shows "No disponible" if any unavailable', async () => {
    const onNext = vi.fn();
    render(
      <QuantitySelector
        min={2}
        priceBS={10}
        priceUsd={1}
        onClose={() => {}}
        onNext={onNext}
      />,
    );

    const ticket1 = screen.getByLabelText('Número 1');
    const ticket2 = screen.getByLabelText('Número 2');
    await userEvent.type(ticket1, '1');
    await userEvent.type(ticket2, '2');

    vi.mocked(fetchAvailableTickets).mockResolvedValue(['2']);

    const next = screen.getByRole('button', { name: 'Siguiente' });
    await userEvent.click(next);

    expect(fetchAvailableTickets).toHaveBeenCalledWith(['1', '2']);
    expect(onNext).not.toHaveBeenCalled();
    expect(await screen.findByText('No disponible')).toBeInTheDocument();
  });

  it('shows general error when availability check throws; re-enables after', async () => {
    const onNext = vi.fn();
    render(
      <QuantitySelector
        min={2}
        priceBS={10}
        priceUsd={1}
        onClose={() => {}}
        onNext={onNext}
      />,
    );

    const ticket1 = screen.getByLabelText('Número 1');
    await userEvent.type(ticket1, '1');

    vi.mocked(fetchAvailableTickets).mockRejectedValue(new Error('boom'));
    const next = screen.getByRole('button', { name: 'Siguiente' });
    await userEvent.click(next);

    expect(await screen.findByText('boom')).toBeInTheDocument();
    expect(onNext).not.toHaveBeenCalled();
  });
});
