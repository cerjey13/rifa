import React, { useEffect, useState } from 'react';
import { ProgressBar } from '@src/landing/components/ProgressBar';

export const CountdownSection = () => {
  const targetDate = new Date();
  targetDate.setHours(targetDate.getHours() + 24);
  const salesPercentage = 77.95;

  return (
    <div className='p-4 space-y-4 max-w-lg mx-auto'>
      <div className='flex flex-wrap items-center gap-3 text-brandLightGray'>
        <span className='flex items-center gap-1'>
          🕑 Fecha del sorteo: ...
        </span>
      </div>

      <div className='bg-brandDark rounded p-3 text-center text-white text-sm font-semibold'>
        Lanzado: 10/7/2025 10:25 PM
      </div>
      <div className='bg-brandDark rounded p-3 text-center text-white text-sm font-semibold'>
        Termina En: 07:18:26
      </div>

      <CountdownTimer targetDate={targetDate} />

      <ProgressBar percentage={salesPercentage} />
      <p className='text-right text-brandLightGray'>
        {salesPercentage}% vendido
      </p>
    </div>
  );
};

interface CountdownTimerProps {
  targetDate: Date;
}

interface LotTimer {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
}

const calculateTimeLeft = (targetDate: Date) => {
  const difference = +targetDate - +new Date();
  let timeLeft: LotTimer = {
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0,
  };

  if (difference > 0) {
    timeLeft = {
      days: Math.floor(difference / (1000 * 60 * 60 * 24)),
      hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
      minutes: Math.floor((difference / 1000 / 60) % 60),
      seconds: Math.floor((difference / 1000) % 60),
    };
  }
  return timeLeft;
};

const CountdownTimer: React.FC<CountdownTimerProps> = ({ targetDate }) => {
  const [timeLeft, setTimeLeft] = useState<LotTimer>(
    calculateTimeLeft(targetDate),
  );
  const units = ['days', 'hours', 'minutes', 'seconds'] as const;
  const labelsEs = {
    days: 'DÍAS',
    hours: 'HORAS',
    minutes: 'MINUTOS',
    seconds: 'SEGUNDOS',
  } as const;

  useEffect(() => {
    const timer = setTimeout(() => {
      setTimeLeft(calculateTimeLeft(targetDate));
    }, 1000);
    return () => clearTimeout(timer);
  }, [timeLeft, targetDate]);

  return (
    <div className='flex justify-center gap-4 text-center font-mono text-base sm:text-xl'>
      {units.map((unit) => (
        <div key={unit} className='flex flex-col'>
          <span className='text-2xl sm:text-4xl font-bold'>
            {timeLeft[unit as keyof typeof timeLeft]}
          </span>
          <span className='uppercase text-xs sm:text-sm'>{labelsEs[unit]}</span>
        </div>
      ))}
    </div>
  );
};
