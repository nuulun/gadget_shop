/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Soft IT Palette
        background: "#F8FAFC", // Very light blue-grey (Soft)
        surface: "#FFFFFF",    // Pure white for cards
        primary: "#3B82F6",    // Tech Blue
        secondary: "#64748B",  // Muted Text (Slate)
        accent: "#8B5CF6",     // Violet (for gradients)
      },
      boxShadow: {
        // Custom shadows for that modern depth
        'soft': '0 4px 20px -2px rgba(0, 0, 0, 0.05)', 
        'glow': '0 0 15px rgba(59, 130, 246, 0.5)',    
      },
      fontFamily: {
        // Ensuring the clean, tech look
        sans: ['Inter', 'sans-serif'], 
      }
    },
  },
  plugins: [],
}