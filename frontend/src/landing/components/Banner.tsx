export const Banner = () => (
  <div
    className='relative rounded-md overflow-hidden my-6 max-w-md mx-auto shadow-lg'
    style={{ backgroundImage: "url('/banner-image.jpg')" }}
  >
    <div className='absolute inset-0 bg-black bg-opacity-50 flex flex-col justify-end p-4'>
      <h2 className='text-white text-xl sm:text-3xl font-bold text-center'>
        Sorteos justos y transparentes para todos
      </h2>
      <div className='flex justify-between mt-4 text-xs text-brandLightGray font-semibold'>
        <span>Resultados el 11 de Julio a las 10:00PM</span>
        <span>Resultado por Lotería del Táchira</span>
        <span>CONALOT</span>
      </div>
    </div>
  </div>
);
