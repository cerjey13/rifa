import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import React from 'react';
import { AppRouter } from './AppRouter';

// Stub children so the router test stays lightweight
vi.mock('@src/landing/LandingPage', () => ({
  LandingPage: () => <div data-testid='landing'>landing</div>,
}));

vi.mock('@src/admin/Admin', () => ({
  AdminApp: () => <div data-testid='admin'>admin</div>,
}));

// Pass-through guard
vi.mock('@src/context/AuthGuard', () => ({
  AuthGuard: ({ children }: { children: React.ReactNode }) => (
    <div data-testid='guard'>{children}</div>
  ),
}));

describe('AppRouter', () => {
  // reset URL before each test
  beforeEach(() => {
    window.history.pushState({}, '', '/');
  });

  it('renders LandingPage at root "/"', async () => {
    window.history.pushState({}, '', '/');
    render(<AppRouter />);

    expect(await screen.findByTestId('landing')).toBeInTheDocument();
  });

  it('wraps AdminApp with AuthGuard at "/dashboard"', async () => {
    window.history.pushState({}, '', '/dashboard');
    render(<AppRouter />);

    // Guard is present and admin page is rendered
    expect(await screen.findByTestId('guard')).toBeInTheDocument();
    expect(screen.getByTestId('admin')).toBeInTheDocument();
  });

  it('redirects unknown routes to "/"', async () => {
    window.history.pushState({}, '', '/nope');
    render(<AppRouter />);

    // After Navigate replace, we should see landing and path reset to "/"
    expect(await screen.findByTestId('landing')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/');
  });
});
