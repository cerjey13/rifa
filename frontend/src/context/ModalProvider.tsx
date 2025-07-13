import { type ReactNode, useState } from 'react';
import { ModalContext } from './useModal';

export interface ModalContextType {
  modalOpen: boolean;
  isRegister: boolean;
  openLoginModal: (register?: boolean) => void;
  closeModal: () => void;
}

export const ModalProvider = ({ children }: { children: ReactNode }) => {
  const [modalOpen, setModalOpen] = useState<boolean>(false);
  const [isRegister, setIsRegister] = useState<boolean>(false);

  const openLoginModal = (register = false) => {
    setIsRegister(register);
    setModalOpen(true);
  };

  const closeModal = () => setModalOpen(false);

  return (
    <ModalContext.Provider
      value={{ modalOpen, isRegister, openLoginModal, closeModal }}
    >
      {children}
    </ModalContext.Provider>
  );
};
