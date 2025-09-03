import { describe, it, expect } from 'vitest';
import { safeArray } from './arrays';
import { formatDateES } from './dates';
import { getErrorMessage } from './errors';

describe('safeArray', () => {
  it('returns the same array when input is an array', () => {
    const input = [1, 2, 3];
    const result = safeArray<number>(input);
    expect(result).toBe(input);
    expect(result).toEqual([1, 2, 3]);
  });

  it('returns [] for non-array values', () => {
    expect(safeArray<number>(null as unknown)).toEqual([]);
    expect(safeArray<number>(undefined as unknown)).toEqual([]);
    expect(safeArray<number>(123 as unknown)).toEqual([]);
    expect(safeArray<number>('abc' as unknown)).toEqual([]);
    expect(safeArray<number>({} as unknown)).toEqual([]);
  });
});

describe('formatDateES', () => {
  it('formats a PM time in es-AR using America/Caracas', () => {
    // 2025-03-15T18:30:00Z -> Caracas = 14:30 (2:30 PM) on 15/03/2025
    const out = formatDateES('2025-03-15T18:30:00.000Z');

    // date part
    expect(out).toMatch(/15\/03\/2025/);

    // time part: allow either "02:30" or "2:30" (some engines)
    expect(out).toMatch(/\b0?2:30\b/);

    // AM/PM marker (handle spaces / NBSP and optional dots): "p. m." / "pm"
    expect(out.toLowerCase()).toMatch(/\bp\.?[\s\u00a0]?m\.?\b/);
  });

  it('formats an AM time in es-AR using America/Caracas', () => {
    // 2025-03-15T04:05:00Z -> Caracas = 00:05 (12:05 AM) on 15/03/2025
    const out = formatDateES('2025-03-15T04:05:00.000Z');

    expect(out).toMatch(/15\/03\/2025/);
    expect(out).toMatch(/\b12:05\b/);
    expect(out.toLowerCase()).toMatch(/\ba\.?[\s\u00a0]?m\.?\b/);
  });

  it('throws on invalid date input', () => {
    //throws RangeError
    expect(() => formatDateES('not-a-date')).toThrow(RangeError);
  });
});

describe('getErrorMessage', () => {
  it('returns message from Error instance', () => {
    const err = new Error('boom');
    expect(getErrorMessage(err, 'fallback')).toBe('boom');
  });

  it('returns the string when error is a string', () => {
    expect(getErrorMessage('bad stuff', 'fallback')).toBe('bad stuff');
  });

  it('returns fallback for non-Error non-string values', () => {
    expect(getErrorMessage(123 as unknown, 'fallback')).toBe('fallback');
    expect(getErrorMessage({} as unknown, 'fallback')).toBe('fallback');
    expect(getErrorMessage(null as unknown, 'fallback')).toBe('fallback');
    expect(getErrorMessage(undefined as unknown, 'fallback')).toBe('fallback');
  });
});
