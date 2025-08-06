import { useState } from 'react';
import { FaInstagram, FaTiktok } from 'react-icons/fa';
import { TermsModal } from '@src/components/Terms/TermsModal';
import { PrivacyModal } from '@src/components/Privacy/PrivacyModal';

export const Footer = () => {
  const [termsOpen, setTermsOpen] = useState(false);
  const [privacyOpen, setPrivacyOpen] = useState(false);

  return (
    <>
      <footer className='bg-[#121726] text-white text-sm px-5 pb-5 mt-5'>
        <div className='max-w-screen-xl mx-auto grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8'>
          <div>
            <img src='/logo.svg' alt='Logo' className='h-10 aspect-[1.22]' />
            <p className='text-gray-400 mt-2'>
              Tu plataforma de confianza para rifas emocionantes. Premios
              increíbles, transparencia total y la mejor experiencia de usuario.
            </p>
          </div>

          <div>
            <h2 className='font-semibold text-white mb-2'>Secciones</h2>
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

          <div>
            <h2 className='font-semibold text-white mb-2'>Contacto</h2>
            <a
              href='https://wa.me/+584141551801'
              target='_blank'
              rel='noopener noreferrer'
              className='text-white underline hover:text-pink-400 focus:outline focus:outline-pink-300'
            >
              04141551801
            </a>
            <p className='text-gray-300'>Caracas, Venezuela</p>
          </div>

          <div>
            <h2 className='font-semibold text-white mb-2'>Síguenos</h2>
            <div className='flex space-x-4 text-xl'>
              <a
                href='https://www.instagram.com/suerteconsarah'
                target='_blank'
                rel='noopener noreferrer'
                aria-label='Instagram de Suerte Con Sarah'
                title='Instagram de Suerte Con Sarah'
                className='text-gray-300 hover:text-pink-500'
              >
                <FaInstagram aria-hidden='true' />
              </a>
              <a
                href='https://www.tiktok.com/@suerteconsarah'
                target='_blank'
                rel='noopener noreferrer'
                className='text-gray-300 hover:text-pink-500'
                aria-label='TikTok de Suerte Con Sarah'
                title='TikTok de Suerte Con Sarah'
              >
                <FaTiktok aria-hidden='true' />
              </a>
            </div>
          </div>
        </div>

        <div className='border-t border-gray-700 mt-10 pt-4 text-center text-xs text-gray-400'>
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
