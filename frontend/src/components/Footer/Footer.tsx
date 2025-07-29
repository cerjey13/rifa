import logo from '@src/assets/logo.jpg';
import { useState } from 'react';
import { FaInstagram, FaTiktok } from 'react-icons/fa';
import { TermsModal } from '../Terms/TermsModal';
import { PrivacyModal } from '../Privacy/PrivacyModal';

export const Footer = () => {
  const [termsOpen, setTermsOpen] = useState(false);
  const [privacyOpen, setPrivacyOpen] = useState(false);

  return (
    <>
      <footer className='bg-[#121726] text-white text-sm px-5 pb-5 mt-5'>
        <div className='max-w-screen-xl mx-auto grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8'>
          {/* Brand/Description */}
          <div>
            <img src={logo} alt='Logo' className='h-10 mb-4' />
            <p className='text-gray-400'>
              Tu plataforma de confianza para rifas emocionantes. Premios
              increíbles, transparencia total y la mejor experiencia de usuario.
            </p>
          </div>

          {/* Secciones */}
          <div>
            <h3 className='font-semibold text-white mb-2'>Secciones</h3>
            <ul className='space-y-1 text-gray-300'>
              <li>
                <button
                  onClick={() => setTermsOpen(true)}
                  className='hover:text-orange-400'
                >
                  Términos y Condiciones
                </button>
              </li>
              <li>
                <button
                  onClick={() => setPrivacyOpen(true)}
                  className='hover:text-orange-400'
                >
                  Política de Privacidad
                </button>
              </li>
            </ul>
          </div>

          {/* Contacto */}
          <div>
            <h3 className='font-semibold text-white mb-2'>Contacto</h3>
            <a
              href='https://wa.me/+584141551801'
              target='_blank'
              rel='noopener noreferrer'
              className='text-gray-300 hover:text-pink-500'
            >
              +584123456789
            </a>
            <p className='text-gray-300'>Caracas, Venezuela</p>
          </div>

          {/* Redes Sociales */}
          <div>
            <h3 className='font-semibold text-white mb-2'>Síguenos</h3>
            <div className='flex space-x-4 text-xl'>
              <a
                href='https://instagram.com'
                target='_blank'
                rel='noopener noreferrer'
                className='text-gray-300 hover:text-pink-500'
              >
                <FaInstagram />
              </a>
              <a
                href='https://tiktok.com'
                target='_blank'
                rel='noopener noreferrer'
                className='text-gray-300 hover:text-white'
              >
                <FaTiktok />
              </a>
            </div>
          </div>
        </div>

        <div className='border-t border-gray-700 mt-10 pt-4 text-center text-xs text-gray-500'>
          © 2025 Rifa Lottery. Todos los derechos reservados. Plataforma segura
          y confiable.
        </div>
      </footer>

      <TermsModal isOpen={termsOpen} onClose={() => setTermsOpen(false)} />
      <PrivacyModal
        isOpen={privacyOpen}
        onClose={() => setPrivacyOpen(false)}
      />
    </>
  );
};
