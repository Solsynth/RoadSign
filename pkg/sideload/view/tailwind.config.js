/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["light"],
          primary: "#4750a3",
          secondary: "#93c5fd",
          accent: "#0f766e",
          info: "#67e8f9",
          success: "#15803d",
          warning: "#f97316",
          error: "#dc2626",
          "--rounded-box": "0",
          "--rounded-btn": "0",
          "--rounded-badge": "0",
          "--tab-radius": "0",
        },
      },
      {
        dark: {
          ...require("daisyui/src/theming/themes")["dark"],
          primary: "#4750a3",
          secondary: "#93c5fd",
          accent: "#0f766e",
          info: "#67e8f9",
          success: "#15803d",
          warning: "#f97316",
          error: "#dc2626",
          "--rounded-box": "0",
          "--rounded-btn": "0",
          "--rounded-badge": "0",
          "--tab-radius": "0",
        },
      },
    ],
  },
  plugins: [require("daisyui")],
};
