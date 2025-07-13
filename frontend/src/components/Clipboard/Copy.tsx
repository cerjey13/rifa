import { useState } from 'react';

interface CopyableTextProps {
  text: string;
  copyText?: string;
}

export const CopyableText = ({ text, copyText }: CopyableTextProps) => {
  const [copied, setCopied] = useState(false);
  const handleCopy = () => {
    navigator.clipboard.writeText(copyText ?? text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <span className='inline-flex items-center space-x-2'>
      <span>{text}</span>
      <button
        onClick={handleCopy}
        aria-label={`Copiar ${copyText ?? text}`}
        className='text-yellow-500 hover:text-yellow-400 focus:outline-none p-2 rounded'
        type='button'
      >
        📋
      </button>
      {copied && <span className='text-green-400 text-sm'>¡Copiado!</span>}
    </span>
  );
};
