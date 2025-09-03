import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';

import { BuyForm } from './BuyForm';
import * as purchaseApi from '@src/api/purchase';

const toastError = vi.fn<(msg: string, opts?: { duration?: number }) => void>();
vi.mock('sonner', () => ({
  toast: { error: (m: string, o?: { duration?: number }) => toastError(m, o) },
}));

vi.mock('@src/components/Clipboard/Copy', () => ({
  CopyableText: ({ text }: { text: string }) => <span>{text}</span>,
}));

const makeFile = (bytes: number, name = 'file.png', type = 'image/png') =>
  new File([new Uint8Array(bytes)], name, { type });

function renderBuyFormUI(
  overrides?: Partial<React.ComponentProps<typeof BuyForm>>,
) {
  const props: React.ComponentProps<typeof BuyForm> = {
    quantity: 2,
    montoBs: 200,
    montoUSD: 20,
    paymentMethod: 'pago movil',
    selectedNumbers: ['1', '42'],
    onBack: vi.fn(),
    onClose: vi.fn(),
    ...overrides,
  };
  const utils = render(<BuyForm {...props} />);
  return { ...utils, props };
}

describe('<BuyForm />', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders Pago Móvil details and selected numbers; shows BS amount', () => {
    const { props } = renderBuyFormUI({
      paymentMethod: 'pago movil',
      selectedNumbers: ['7', '9'],
      montoBs: 350,
    });

    expect(screen.getByText(/Metodo:/i)).toBeInTheDocument();
    expect(screen.getByText(/PAGO MOVIL/i)).toBeInTheDocument();

    expect(screen.getByText(/BANESCO 0134/)).toBeInTheDocument();
    expect(screen.getByText(/30606459/)).toBeInTheDocument();
    expect(screen.getByText(/04141551801/)).toBeInTheDocument();

    expect(screen.getByText(/Monto BS:\s*350/)).toBeInTheDocument();

    expect(screen.getByText(/Números seleccionados:/)).toBeInTheDocument();
    expect(screen.getByText(/7,\s*9/)).toBeInTheDocument();

    expect(
      screen.getByText(new RegExp(`Cantidad de tickets:\\s*${props.quantity}`)),
    ).toBeInTheDocument();
  });

  it('renders Zelle details and shows USD amount', () => {
    renderBuyFormUI({ paymentMethod: 'zelle', montoUSD: 25 });

    expect(screen.getByText(/Metodo:/i)).toBeInTheDocument();
    expect(screen.getByText(/ZELLE/i)).toBeInTheDocument();

    expect(screen.getByText(/3802389306/)).toBeInTheDocument();
    expect(screen.getByText(/Vicente Méndez/)).toBeInTheDocument();

    expect(screen.getByText(/Monto \(\$\):\s*25/)).toBeInTheDocument();
  });

  it('digits input accepts only digits and trims to max length of 6', async () => {
    renderBuyFormUI();

    const digits = screen.getByLabelText(
      'Últimos 6 Dígitos (Transacción)',
    ) as HTMLInputElement;
    await userEvent.type(digits, '12ab34!@#56xyz');
    expect(digits.value).toBe('123456');
  });

  it('rejects files over 3MB and shows toast error (no file name shown)', async () => {
    const { container } = renderBuyFormUI();
    const fileInput = container.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;

    await userEvent.upload(fileInput, makeFile(3 * 1024 * 1024 + 1, 'big.png'));
    expect(toastError).toHaveBeenCalledWith(
      'El archivo no debe superar los 3 MB.',
      { duration: 5000 },
    );
    expect(
      screen.queryByText(/Archivo seleccionado:\s*big\.png/),
    ).not.toBeInTheDocument();
  });

  it('blocks submit when digits are missing (native validation) and does not call API', async () => {
    vi.spyOn(purchaseApi, 'submitPurchase').mockResolvedValueOnce(undefined);
    renderBuyFormUI();

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await userEvent.upload(fileInput, makeFile(1024, 'ok.png'));

    await userEvent.click(screen.getByRole('button', { name: 'Finalizar' }));

    const digits = screen.getByLabelText(
      'Últimos 6 Dígitos (Transacción)',
    ) as HTMLInputElement;
    expect(digits).toBeInvalid();
    expect(purchaseApi.submitPurchase).not.toHaveBeenCalled();
  });

  it('blocks submit when file is missing (native validation) and does not call API', async () => {
    vi.spyOn(purchaseApi, 'submitPurchase').mockResolvedValueOnce(undefined);

    renderBuyFormUI();

    const digits = screen.getByLabelText(
      'Últimos 6 Dígitos (Transacción)',
    ) as HTMLInputElement;
    await userEvent.type(digits, '123456');

    await userEvent.click(screen.getByRole('button', { name: 'Finalizar' }));

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    expect(fileInput).toBeInvalid();
    expect(purchaseApi.submitPurchase).not.toHaveBeenCalled();
  });

  it('submits, shows success overlay, and Close calls onClose', async () => {
    vi.spyOn(purchaseApi, 'submitPurchase').mockResolvedValueOnce(undefined);

    const { container, props } = renderBuyFormUI({
      paymentMethod: 'pago movil',
    });
    // disable validation for test only
    (container.querySelector('form') as HTMLFormElement).noValidate = true;

    const digits = screen.getByLabelText(
      'Últimos 6 Dígitos (Transacción)',
    ) as HTMLInputElement;
    await userEvent.type(digits, '123456');
    expect(digits.value).toBe('123456');

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await userEvent.upload(fileInput, makeFile(2048, 'ok.png'));

    await userEvent.click(screen.getByRole('button', { name: 'Finalizar' }));

    await waitFor(() => expect(purchaseApi.submitPurchase).toHaveBeenCalled());
    const [payload] = vi.mocked(purchaseApi.submitPurchase).mock.calls[0];

    expect(payload.quantity).toBe(2);
    expect(payload.montoBs).toBe('200.00');
    expect(payload.montoUSD).toBe('20.00');
    expect(payload.paymentMethod).toBe('pago movil');
    expect(payload.transactionDigits).toBe('123456');
    expect(payload.selectedNumbers).toEqual(['1', '42']);
    expect(payload.paymentScreenshot).toBeInstanceOf(File);
    expect(payload.paymentScreenshot.name).toBe('ok.png');

    // Success overlay visible (aria-hidden should be false) and has "¡Éxito!"
    const successHeading = await screen.findByText('¡Éxito!');
    const overlay = successHeading.closest('[aria-hidden]');
    expect(overlay).toHaveAttribute('aria-hidden', 'false');

    // Close button calls onClose
    const closeBtn = screen.getByRole('button', { name: 'Cerrar' });
    await userEvent.click(closeBtn);
    expect(props.onClose).toHaveBeenCalledTimes(1);
  });

  it('shows error message when submit fails and resets loading', async () => {
    vi.spyOn(purchaseApi, 'submitPurchase').mockRejectedValueOnce(
      new Error('nope'),
    );
    const { container } = renderBuyFormUI({
      paymentMethod: 'pago movil',
    });
    // disable validation for test only
    (container.querySelector('form') as HTMLFormElement).noValidate = true;

    const digits = screen.getByLabelText(
      'Últimos 6 Dígitos (Transacción)',
    ) as HTMLInputElement;
    await userEvent.type(digits, '123456');

    const fileInput = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await userEvent.upload(fileInput, makeFile(2048, 'ok.png'));

    await userEvent.click(screen.getByRole('button', { name: 'Finalizar' }));

    await waitFor(() =>
      expect(vi.mocked(purchaseApi.submitPurchase)).toHaveBeenCalledTimes(1),
    );

    expect(
      await screen.findByText('Error al enviar los datos, intenta nuevamente'),
    ).toBeInTheDocument();
    // Button text should be back to "Finalizar" (no spinner)
    expect(
      screen.getByRole('button', { name: 'Finalizar' }),
    ).toBeInTheDocument();
  });

  it('shows loading spinner and disables submit while in-flight', async () => {
    let resolve: (() => void) | undefined;
    const deferred = new Promise<void>((res) => {
      resolve = res;
    });

    vi.spyOn(purchaseApi, 'submitPurchase').mockReturnValueOnce(deferred);

    const { container } = renderBuyFormUI({ paymentMethod: 'pago movil' });

    (container.querySelector('form') as HTMLFormElement).noValidate = true;

    await userEvent.type(
      screen.getByLabelText('Últimos 6 Dígitos (Transacción)'),
      '123456',
    );
    const fileInput = container.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    await userEvent.upload(fileInput, makeFile(1024, 'ok.png'));

    const submitBtn = screen.getByRole('button', { name: 'Finalizar' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
      expect(submitBtn).toHaveAttribute('aria-busy', 'true');
      expect(submitBtn.querySelector('svg.animate-spin')).toBeTruthy();
    });

    resolve?.();
    // After resolve, success overlay should appear
    await screen.findByText('¡Éxito!');
  });

  it('Back button calls onBack', async () => {
    const onBack = vi.fn();
    renderBuyFormUI({ onBack });

    await userEvent.click(screen.getByRole('button', { name: 'Atrás' }));
    expect(onBack).toHaveBeenCalledTimes(1);
  });
});
