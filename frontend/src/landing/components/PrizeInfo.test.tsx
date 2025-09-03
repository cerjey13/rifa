import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { PrizeInfo } from './PrizeInfo';

vi.mock('@src/assets/premio-500.webp', () => ({
  default: 'premio-500.webp',
}));

describe('<PrizeInfo />', () => {
  it('renders title, subtitle and image', () => {
    render(<PrizeInfo />);
    expect(
      screen.getByRole('heading', { name: /Prueba tu suerte/i }),
    ).toBeInTheDocument();
    expect(screen.getByText(/Participa y gana/i)).toBeInTheDocument();

    const img = screen.getByAltText('rifa') as HTMLImageElement;
    expect(img).toBeInTheDocument();
    expect(img.src).toContain('premio-500.webp');
  });
});
