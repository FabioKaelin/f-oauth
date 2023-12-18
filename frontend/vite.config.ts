import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslintPlugin from "vite-plugin-eslint"

const path = require("path")
export default defineConfig({
    plugins: [vue(), eslintPlugin()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src")
        },
        extensions: [".mjs", ".js", ".ts", ".jsx", ".tsx", ".json", ".vue"]
    },
    build: {
        rollupOptions: {
            output: {
                manualChunks(id) {
                    if (id.includes("node_modules")) {
                        return id.toString().split("node_modules/")[1].split("/")[0].toString()
                    }
                }
            }
            // input: {
            //     main: path.resolve(__dirname, "index.html"),
            //     // RulesView.vue should be included in the bundle
            //     rules: path.resolve(__dirname, "src/views/RulesView.vue"),
            // }
        }
        // ssr: true,
    },
    server: {
        port: 5173,
        open: true
    },
    preview: {
        port: 5173,
        host: "0.0.0.0",
        cors: true
        // proxy: {
        //     "/api": {
        //         target: "http://localhost:3000",
        //         changeOrigin: true,
        //         rewrite: (path) => path.replace(/^\/api/, "")
        //     }
        // }
    }
    // ssr: {
    //     noExternal: ["vue", "vue-router"]
    // },
    // css: {
    //     preprocessorOptions: {
    //         scss: {
    //             additionalData: `@import "@/assets/css/variables.scss";`
    //         },
    //     }
    // }
})
