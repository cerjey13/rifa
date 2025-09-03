import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import CountdownSection from './CountdownSection';

vi.mock('@src/landing/components/ProgressBar', () => ({
  ProgressBar: ({ percentage }: { percentage: number }) => (
    <div data-testid='bar' data-percentage={String(percentage)} />
  ),
}));

describe('<CountdownSection />', () => {
  it('shows loading message', () => {
    render(
      <CountdownSection percentage={12.34} loading={true} error={false} />,
    );
    expect(screen.getByText(/Cargando disponibilidad/i)).toBeInTheDocument();
  });

  it('shows 0% when error=true and still renders progress bar with original percentage', () => {
    render(<CountdownSection percentage={55.5} loading={false} error={true} />);
    expect(screen.getByTestId('bar').dataset.percentage).toBe('55.5');
    expect(screen.getByText('0% Vendido')).toBeInTheDocument();
  });

  it('renders percentage with two decimals when not error', () => {
    render(
      <CountdownSection percentage={12.3} loading={false} error={false} />,
    );
    expect(screen.getByText('12.30% Vendido')).toBeInTheDocument();
    expect(screen.getByTestId('bar').dataset.percentage).toBe('12.3');
  });
});
