/// <reference types="vite/client" />
interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_MONTO_BS: string;
  readonly VITE_MONTO_USD: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
