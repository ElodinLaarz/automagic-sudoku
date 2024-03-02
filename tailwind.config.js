const colors = require('tailwindcss/colors')
const fonts = require('tailwindcss/defaultTheme').fontFamily

module.exports = {
  mode: 'jit',
  content: [
    './src/**/*.{go,html,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: colors.rose,
        secondary: colors.slate,
      },
      fonts: {
        default: ['Open Sans', ...fonts.sans],
      }
    },
  },
  plugins: [],
}

