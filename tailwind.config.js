// tailwind.config.js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/**/*.{html,js,templ}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        dark: {
          base: '#171717',      // Very dark gray for background
          paper: '#262626',     // Slightly lighter for cards
          accent: '#404040',    // For borders and dividers
          text: {
            primary: '#F5F5F5', // Almost white for primary text
            secondary: '#A3A3A3' // Gray for secondary text
          }
        },
        pastel: {
          blue: '#AED9E0',    // Soft blue
          green: '#B8E0D2',   // Mint green
          pink: '#FAD2E1',    // Soft pink
          yellow: '#FBE5C8',  // Warm yellow
          purple: '#D4C4E9',  // Gentle purple
          background: '#FAFBFF', // Very light blue-white
          // New eye-friendly colors
          base: '#F5F3EF',     // Very soft warm white (main background)
          paper: '#FAF9F6',    // Slightly warmer white for cards
          warmGray: '#EFEAE4', // Warm gray for sections
          text: '#4A4743',     // Soft brown-gray for text
        },
        primary: {
          50: '#f0f9ff',  // Pastel blue
          100: '#e0f2fe',
          200: '#bae6fd',
          300: '#7dd3fc',
          400: '#38bdf8',
          500: '#0ea5e9',
        },
        secondary: {
          50: '#fdf4ff',  // Pastel purple
          100: '#fae8ff',
          200: '#f5d0fe',
          300: '#f0abfc',
          400: '#e879f9',
          500: '#d946ef',
        },
        accent: {
          50: '#fff1f2',  // Pastel pink
          100: '#ffe4e6',
          200: '#fecdd3',
          300: '#fda4af',
          400: '#fb7185',
          500: '#f43f5e',
        },
        neutral: {
          50: '#fafafa',
          100: '#f5f5f5',
          200: '#e5e5e5',
          300: '#d4d4d4',
          400: '#a3a3a3',
          500: '#737373',
          600: '#525252',
          700: '#404040',
          800: '#262626',
          900: '#171717',
        }
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'monospace'],
        sans: ['Inter', 'sans-serif'],
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: '65ch',
            color: '#404040',
            a: {
              color: '#0ea5e9',
              '&:hover': {
                color: '#38bdf8',
              },
            },
          },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
  darkMode: 'class',
}
