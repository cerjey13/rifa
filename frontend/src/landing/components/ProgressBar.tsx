import React from 'react';

interface ProgressBarProps {
  percentage: number; // 0-100
}

export const ProgressBar: React.FC<ProgressBarProps> = ({ percentage }) => (
  <div className='w-full max-w-md mx-auto bg-gray-300 rounded-full h-4 overflow-hidden my-4'>
    <div
      className='bg-orange-500 h-4 rounded-full transition-all duration-500'
      style={{ width: `${percentage ?? 0}%` }}
    />
  </div>
);
