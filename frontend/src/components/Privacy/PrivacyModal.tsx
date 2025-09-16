import { useEffect } from 'react';

interface PrivacyModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const PrivacyModal = ({ isOpen, onClose }: PrivacyModalProps) => {
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    if (isOpen) window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const privacyItems = [
    {
      title: 'Información que recopilamos',
      content:
        'Recopilamos los siguientes datos cuando interactúas con nuestro sitio web o servicios: \n- Información personal: nombre, correo electrónico, número telefónico y otros datos que proporciones voluntariamente al registrarte o usar nuestros servicios.',
    },
    {
      title: 'Uso de la información',
      content:
        'Utilizamos tus datos para: \n- Enviarte notificaciones, promociones y actualizaciones relevantes. \n- Personalizar y mejorar tu experiencia en nuestra plataforma. \n- Cumplir con obligaciones legales o resolver disputas.',
    },
    {
      title: 'Compartir y divulgar la información',
      content:
        'No vendemos ni alquilamos tu información. Solo la compartimos con: \n- Proveedores de servicios como plataformas de email, análisis de datos o soporte técnico, bajo acuerdos de confidencialidad. \n- Autoridades legales cuando sea requerido por ley o en cumplimiento de procesos judiciales.',
    },
    {
      title: 'Seguridad de la información',
      content:
        'Implementamos medidas de seguridad razonables para proteger tus datos contra accesos no autorizados, pérdidas o alteraciones. Sin embargo, ningún sistema es 100% seguro. Nos esforzamos al máximo por mantener la integridad y confidencialidad de tu información.',
    },
    {
      title: 'Tus derechos',
      content:
        'Puedes: \n- Acceder a los datos que tenemos sobre ti. \n- Solicitar su corrección o eliminación (salvo excepciones legales). \n- Oponerte a su procesamiento en ciertas circunstancias.   ',
    },
    {
      title: 'Cookies',
      content:
        'Este sitio utiliza cookies para mejorar tu experiencia y analizar el uso del sitio. Puedes configurar tu navegador para rechazarlas o notificarte al recibirlas.',
    },
    {
      title: 'Cambios en la política',
      content:
        'Nos reservamos el derecho de modificar esta política en cualquier momento. Las actualizaciones estarán disponibles en esta página y serán notificadas si es necesario. Te recomendamos revisarla periódicamente.',
    },
  ];

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70 p-4'>
      <div className='bg-[#111827] text-white p-6 rounded-lg w-full max-w-3xl max-h-[90vh] overflow-y-auto shadow-xl relative'>
        {/* Close Button */}
        <button
          onClick={onClose}
          className='absolute top-4 right-4 text-gray-400 hover:text-white text-2xl'
        >
          &times;
        </button>

        {/* Header */}
        <div className='pl-4 mb-6'>
          <h2 className='text-2xl font-bold text-[#facc15]'>
            Política de Privacidad
          </h2>
        </div>

        {/* Content */}
        <div className='space-y-6 text-sm text-gray-300 leading-relaxed pr-1'>
          <p>
            En Suerte con Sarah nos comprometemos a proteger tu privacidad y
            garantizar que tu experiencia con nuestros servicios sea segura y
            confiable.
          </p>

          {privacyItems.map((item, index) => (
            <div key={index}>
              <h3 className='text-[#facc15] font-semibold'>
                {index + 1}. {item.title}
              </h3>
              <p className='mt-1 whitespace-pre-line'>{item.content}</p>
            </div>
          ))}
        </div>

        {/* Footer Button */}
        <div className='mt-8 flex justify-end'>
          <button
            onClick={onClose}
            className='bg-[#2f00ff] px-5 py-2 rounded-md text-white font-semibold shadow'
          >
            Cerrar
          </button>
        </div>
      </div>
    </div>
  );
};
