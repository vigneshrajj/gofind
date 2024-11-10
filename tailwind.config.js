/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./static/**/*.{js,html}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
};
