/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx,vue}",
    ],
    theme: {
        extend: {},
    },
    plugins: [require("@tailwindcss/typography"), require("daisyui")],
    daisyui: {
        themes: [{
            winter: {
                ...require("daisyui/src/theming/themes")["winter"],
                // primary: "#2dd4bf",
                primary: "#3bad9d",
            }
        }, {
            dim: {
                ...require("daisyui/src/theming/themes")["dim"],
                // primary: "#2dd4bf",
                primary: "#3bad9d",

            }
        }],
        darkTheme: "dim",
        base: true,
        styled: true,
        utils: true,
        prefix: "",
        logs: true,
        themeRoot: ":root"
    }
}

