import { createContext, useContext } from 'react';
import type { ModalContextType } from './ModalProvider';

export const ModalContext = createContext<ModalContextType | undefined>(
  undefined,
);

export const useModal = () => {
  const context = useContext<ModalContextType | undefined>(ModalContext);
  if (!context) throw new Error('useModal must be used within ModalProvider');
  return context;
};
