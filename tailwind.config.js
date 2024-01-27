/** @type {import('tailwindcss').Config} */
export const content = ["./internal/views/**/*.{html,js,templ}"];
export const theme = {
     extend: {
          colors: {
               "blue-gray": {
                    50: "#ECEFF1",
                    900: "#263238",
                    800: "#37474F",
                    700: "#455A64",
                    600: "#546E7A",
                    500: "#607D8B",
                    400: "#78909C",
                    300: "#90A4AE",
                    200: "#B0BEC5",
                    100: "#CFD8DC",
               }
          },
     },
};
export const plugins = [];

