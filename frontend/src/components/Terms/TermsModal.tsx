import { useEffect } from 'react';

interface TermsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const TermsModal = ({ isOpen, onClose }: TermsModalProps) => {
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    if (isOpen) window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const terms = [
    {
      title: 'Registro y Datos Personales',
      content:
        'Nombre, teléfono, correo electrónico. Mantener información actualizada es responsabilidad del usuario.',
    },
    {
      title: 'Comprobantes de Pago',
      content:
        'Ganadores deben presentar comprobante válido para validar su premio.',
    },
    {
      title: 'Compra Mínima de Boletos',
      content:
        'Para efectos de participación válida en el sorteo, el usuario deberá adquirir un mínimo de dos (2) boletos por operación. No se aceptarán compras inferiores a esta cantidad.',
    },
    {
      title: 'Datos Adicionales',
      content: 'Se solicitará orden de compra y captura de pago.',
    },
    {
      title: 'Redes Sociales',
      content:
        'Se solicitará cuenta de Instagram del ganador para verificar su identidad.',
    },
    {
      title: 'Edad',
      content: 'Solo mayores de 18 años pueden participar.',
    },
    {
      title: 'Participantes Internacionales',
      content: 'Requieren representante legal dentro de Venezuela.',
    },
    {
      title: 'Plazo para Retirar',
      content:
        'Máximo 15 días naturales desde el anuncio para reclamar el premio.',
    },
    {
      title: 'Correo Electrónico',
      content: 'Será utilizado para notificaciones, promociones y sorteos.',
    },
    {
      title: 'Cambios de Números',
      content:
        'No se permiten cambios una vez aprobado el número. Debes comprar otro.',
    },
    {
      title: 'Métodos de Pago',
      content:
        'Solo se aceptan los medios oficiales proporcionados en la plataforma.',
    },
    {
      title: 'Lotería',
      content:
        'Usamos los resultados de la Lotería del Táchira (Resultados Super Gana).',
      link: 'https://tripletachira.com',
    },
    {
      title: 'Autorización de Uso de Imagen de los Ganadores',
      content:
        'Al participar y resultar ganador, el usuario acepta ceder su imagen para la creación de contenido audiovisual relacionado con el sorteo, así como su difusión en redes sociales y en la entrega del premio. Esta condición es obligatoria.',
    },
    {
      title: 'Comprobantes Incorrectos',
      content:
        'Si el comprobante es ilegible o incorrecto, se notificará al usuario. Si no se corrige a tiempo, el ticket será anulado y el monto quedará como abono.',
    },
    {
      title: 'Modificación de Términos',
      content:
        'Nos reservamos el derecho de modificar los términos. Los cambios serán notificados.',
    },
    {
      title: 'Política de Reembolso',
      content:
        'No se devuelve dinero. El monto quedará como abono en la plataforma para próximos eventos.',
    },
  ];

  return (
    <div className='fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70 p-4'>
      <div className='bg-[#1c1f2b] text-white p-6 rounded-lg w-full max-w-3xl max-h-[90vh] overflow-y-auto shadow-xl relative'>
        <button
          onClick={onClose}
          className='absolute top-4 right-4 text-gray-400 hover:text-white text-2xl'
        >
          &times;
        </button>

        <div className='pl-4 mb-6'>
          <h2 className='text-2xl font-bold text-[#facc15]'>
            Términos y Condiciones de Uso
          </h2>
        </div>

        <div className='space-y-4 text-sm text-gray-300 leading-relaxed pr-1'>
          <p>
            Bienvenido a nuestra plataforma de rifas. Al acceder y utilizar
            nuestros servicios, aceptas los siguientes términos:
          </p>

          {terms.map((term, index) => (
            <div key={index}>
              <h3 className='text-[#facc15] font-semibold'>
                {index + 1}. {term.title}
              </h3>
              <p className='mt-1'>{term.content}</p>
              {term.link ? (
                <a
                  href='https://tripletachira.com'
                  target='_blank'
                  rel='noopener noreferrer'
                  className='text-orange-400 underline'
                >
                  Tripletachira.com
                </a>
              ) : null}
            </div>
          ))}

          <p className='text-xs text-center text-gray-500 mt-6'>
            © 2025 | Todos los derechos reservados •{' '}
          </p>
        </div>

        {/* Footer Button */}
        <div className='mt-2 flex justify-end'>
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
