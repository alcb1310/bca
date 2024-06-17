/** @type {import('tailwindcss').Config} */
    module.exports = {
        content: [
            "./**/*.html",
            "./**/*.go",
            "./**/*.templ",
        ],
        theme: {
            colors: {
                'blue-gray': {
                    50: '#eceff1',
                    100: '#cfd8dc',
                    200: '#b0bec5',
                    300: '#90a4ae',
                    400: '#78909c',
                    500: '#607d8b',
                    600: '#546e7a',
                    700: '#455a64',
                    800: '#37474f',
                    900: '#263238',
                },
                'green': {
                    500: '#22c55e',
                    600: '#16a34a',
                    700: '#15803d',
                    800: '#166534',
                    900: '#14532d',
                },
                'red' : {
                    500: '#ef4444'
                }
            },
            extend: {
            },
        },
        plugins: [],
    }

