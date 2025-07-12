/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        brandOrange: '#FF7F00',
        brandDark: '#1B1B1F',
        brandGray: '#2E2E38',
        brandLightGray: '#BDBDBD',
        buttonsGray: '#272525',
      },
      fontFamily: {
        sans: ['Poppins', 'sans-serif'],
      },
    },
  },
  plugins: [],
};
