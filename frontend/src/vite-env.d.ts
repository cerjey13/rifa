/// <reference types="vite/client" />
interface ImportMetaEnv {
  readonly API_URL: string;
  readonly montoBs: number;
  readonly montoUsd: number;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
