import premio512 from '@src/assets/premio-500.webp';

export const PrizeInfo = () => (
  <section className='text-center p-4 my-4 space-y-3 max-w-lg mx-auto text-brandLightGray'>
    <h1 className='text-2xl sm:text-3xl font-extrabold leading-tight text-white uppercase'>
      🍀 Prueba tu suerte 🍀
    </h1>
    <p className='text-lg sm:text-xl text-white font-semibold'>
      Participa y gana!
    </p>
    <div className='flex justify-center mx-auto max-w-[512px]'>
      <picture>
        <source
          srcSet={premio512}
          media='(min-width: 641px)'
          type='image/webp'
        />
        <img
          src={premio512}
          alt='rifa'
          width={512}
          height={682}
          loading='eager'
          fetchPriority='high'
          className='w-full max-w-xs sm:max-w-sm rounded-lg shadow-md object-cover'
        />
      </picture>
    </div>
    <p className='text-white font-medium flex justify-center gap-2 items-center'>
      <span>🏆 Premios:</span>
    </p>
    <p className='text-white font-medium flex justify-center gap-2 items-center'>
      <span>🥇 1er Premio: New Outlook II 🏍️ 0 km 2025</span>
    </p>
    <p className='text-white font-medium flex justify-center gap-2 items-center'>
      <span>🥈 2do Premio: El comprador de más boletos 500$💸</span>
    </p>
    <p className='text-white font-medium flex justify-center gap-2 items-center'>
      <span>
        🥉 3er y 🏅4to Premio: Por aproximación (un número antes y uno después
        del número ganador) 100$ 💵
      </span>
    </p>
  </section>
);
