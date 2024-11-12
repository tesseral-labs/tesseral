/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{ts,tsx}', './public/index.html'],
  darkMode: 'class',
  plugins: [require('tailwindcss-animate')],
  prefix: '',
  purge: {
    enabled: true,
    content: ['./src/**/*.{ts,tsx}', './public/index.html'],
  },
  theme: {
    extend: {},
  },
}
