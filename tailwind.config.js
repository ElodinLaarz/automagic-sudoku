const colors = require('tailwindcss/colors')
const fonts = require('tailwindcss/defaultTheme').fontFamily

module.exports = {
  mode: 'jit',
  content: [
    './src/**/*.{go,html,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      fontSize: {
        'variable': 'clamp(2px, 1vw, 1rem)',
      },
      colors: {
        primary: {
          '50': '#fbf6f5',
          '100': '#f6ecea',
          '200': '#f0dcd8',
          '300': '#e4c3bd',
          '400': '#d3a096',
          '500': '#ba7264',
          '600': '#aa6558',
          '700': '#8e5347',
          '800': '#77463d',
          '900': '#643f38',
          '950': '#351e1a',
      },
        secondary: colors.slate,
      },
      fonts: {
        default: ['Open Sans', ...fonts.sans],
      }
    },
  },
  plugins: [],
}

