const colors = require("tailwindcss/colors")

module.exports = {
    purge: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
    darkMode: false, // or 'media' or 'class'
    theme: {
        extend: {},
        colors: {
            green: colors.green,
            black: colors.black,
            gray: colors.gray,
            red: colors.red,
            yellow: colors.yellow
        }
    },
    variants: {
        extend: {},
    },
    plugins: [],
}
