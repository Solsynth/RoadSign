/** @type {import("tailwindcss").Config} */
export default {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {}
  },
  daisyui: {
    themes: [
      {
        light: {
          primary: "#4750a3",
          secondary: "#93c5fd",
          accent: "#0f766e",
          info: "#67e8f9",
          success: "#15803d",
          warning: "#f97316",
          error: "#dc2626",
          neutral: "#2B3440",
          "secondary-content": "oklch(98.71% 0.0106 342.55)",
          "neutral-content": "#D7DDE4",
          "base-100": "oklch(100% 0 0)",
          "base-200": "#F2F2F2",
          "base-300": "#E5E6E6",
          "base-content": "#1f2937",
          "color-scheme": "light",
          "--rounded-box": "0",
          "--rounded-btn": "0",
          "--rounded-badge": "0",
          "--tab-radius": "0"
        }
      },
      {
        dark: {
          primary: "#4750a3",
          secondary: "#93c5fd",
          accent: "#0f766e",
          info: "#67e8f9",
          success: "#15803d",
          warning: "#f97316",
          error: "#dc2626",
          neutral: "#2a323c",
          "neutral-content": "#A6ADBB",
          "base-100": "#1d232a",
          "base-200": "#191e24",
          "base-300": "#15191e",
          "base-content": "#A6ADBB",
          "color-scheme": "dark",
          "--rounded-box": "0",
          "--rounded-btn": "0",
          "--rounded-badge": "0",
          "--tab-radius": "0"
        }
      }
    ]
  },
  plugins: [require("daisyui")]
}
