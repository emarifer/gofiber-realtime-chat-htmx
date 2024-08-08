const { fontFamily } = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */

module.exports = {
    content: [
        "../views/**/*.tmpl"
    ],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Kanit', ...fontFamily.sans]
            }
        },
    },
    plugins: [
        require("daisyui")
    ],
    daisyui: {
        themes: ["dark"]
    }
}

/* TAILWIND COMMANDS EXECUTED FROM THE ROOT FOLDER */

/* npx tailwindcss init web/tailwind/tailwind.config.js */

/* npm i --prefix web/tailwind/ */

/* npm run --prefix web/tailwind/ watch-css */