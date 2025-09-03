import { render } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ProgressBar } from './ProgressBar';

describe('<ProgressBar />', () => {
  it('sets inner width style according to percentage', () => {
    const { container } = render(<ProgressBar percentage={42} />);
    const inner = container.querySelector('.bg-orange-500') as HTMLDivElement;
    expect(inner).toBeTruthy();
    expect(inner).toHaveStyle({ width: '42%' });
  });
});
