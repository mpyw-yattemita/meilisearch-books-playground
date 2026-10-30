import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // The Go API server (backend/cmd/server).
      "/api": "http://localhost:8080",
    },
  },
});
