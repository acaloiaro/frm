
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'ui/**/*.templ',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        /* 
          Below color shades come from https://uicolors.app/create 

          From the PastSight branding kit
          Primary original shade: #00255A
          Secondary original shade: 00A583
        */
        transparent: 'transparent',
        current: 'currentColor',
        'white': '#ffffff',
        'red': {
          '100': '#ffeeee',
          '200': '#ffdddd',
          '300': '#ffbbbb',
          '400': '#ff9999',
          '500': '#ff7777',
          '600': '#ff5555',
          '700': '#ff3333',
          '800': '#ff1111',
          '900': '#ee0000'
        },
        'primary': {
          '50': '#e9f9ff',
          '100': '#cef1ff',
          '200': '#a7e8ff',
          '300': '#6bdeff',
          '400': '#26c6ff',
          '500': '#009fff',
          '600': '#0075ff',
          '700': '#005aff',
          '800': '#004de6',
          '900': '#0047b3',
          '950': '#00255a',
        },
        'secondary': {
          '50': '#ebfef7',
          '100': '#d0fbe9',
          '200': '#a4f6d9',
          '300': '#6aebc5',
          '400': '#2fd8ac',
          '500': '#0abf96',
          '600': '#00a583',
          '700': '#007c66',
          '800': '#036251',
          '900': '#045044',
          '950': '#012d28',
        },
        'gray': {
          '50': '#f6f6f6',
          '100': '#e7e7e7',
          '200': '#d1d1d1',
          '300': '#b0b0b0',
          '400': '#888888',
          '500': '#6d6d6d',
          '600': '#5d5d5d',
          '700': '#4c4c4c',
          '800': '#454545',
          '900': '#3d3d3d',
          '950': '#262626',
        },
      },
      fontFamily: {
        sans: ['Helvetica', 'Arial', 'sans-serif'],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],
  corePlugins: {
    preflight: true,
  }
}

