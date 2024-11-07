/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        './static/tmpl/*.tmpl',
        './static/js/*.js'
    ],
    theme: {
        extend: {
            keyframes: {
                'spin-slow': {
                    '0%': { transform: 'rotate(0deg)' },
                    '100%': { transform: 'rotate(360deg)' },
                },
            },
            animation: {
                'spin-slow': 'spin-slow 3s linear infinite',
            },
            colors: {
                'custom-teal': '#38B2AC',
                'custom-yellow': '#ECC94B',
                'custom-teal-dark': '#2C7A7B',
                'custom-yellow-dark': '#D69E2E',
                'main-bg-start': '#004b49',  // Custom start color for gradient
                'main-bg-end': '#004040'     // Custom end color for gradient
            },
        },
    },
    plugins: [],
}
