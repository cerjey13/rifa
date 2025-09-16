import { afterEach, describe, expect, it, vi } from 'vitest';

describe('config', () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.resetModules();
  });

  it('reads values from import.meta.env and coerces numbers', async () => {
    vi.stubEnv('VITE_API_URL', 'https://api.example.com');
    vi.stubEnv('VITE_MONTO_BS', '350');
    vi.stubEnv('VITE_MONTO_USD', '1');

    vi.resetModules();
    const mod = await import('./config');

    expect(mod.API_URL).toBe('https://api.example.com');
    expect(mod.MONTO_BS).toBe(350);
    expect(mod.MONTO_USD).toBe(1);
  });

  it('handles missing envs: string is undefined, numbers become NaN', async () => {
    vi.stubEnv('VITE_API_URL', undefined);
    vi.stubEnv('VITE_MONTO_BS', 'not-a-num');
    vi.stubEnv('VITE_MONTO_USD', '');

    vi.resetModules();
    const mod = await import('./config');

    expect(mod.API_URL).toBeUndefined();
    expect(Number.isNaN(mod.MONTO_BS)).toBe(true);
    expect(mod.MONTO_USD).toBe(0);
  });
});
