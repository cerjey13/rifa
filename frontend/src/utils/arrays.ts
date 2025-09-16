export function safeArray<T>(data: unknown): T[] {
  return Array.isArray(data) ? data : [];
}
