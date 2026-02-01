/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#2563EB',
        success: '#16A34A',
        warning: '#F59E0B',
        error: '#DC2626',
        surface: {
          page: '#F9FAFB',
          card: '#FFFFFF',
        },
        text: {
          primary: '#111827',
          muted: '#6B7280',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      fontSize: {
        'h1': ['2rem', { lineHeight: '2.5rem', fontWeight: '700' }],
        'h2': ['1.5rem', { lineHeight: '2rem', fontWeight: '600' }],
        'h3': ['1.25rem', { lineHeight: '1.75rem', fontWeight: '600' }],
        'body': ['1rem', { lineHeight: '1.5rem', fontWeight: '400' }],
        'body-sm': ['0.875rem', { lineHeight: '1.25rem', fontWeight: '400' }],
      },
    },
  },
  plugins: [],
}
