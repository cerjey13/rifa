export function getErrorMessage(error: unknown, message: string): string {
  if (error instanceof Error) return error.message;
  if (typeof error === 'string') return error;
  return message;
}
