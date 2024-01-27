import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslintPlugin from "vite-plugin-eslint"


// apple-touch-icon
// <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">

const path = require("path")
// export default defineConfig(() =>({
export default defineConfig(({mode}) =>({
    plugins: [vue(), eslintPlugin()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src")
        },
        extensions: [".mjs", ".js", ".ts", ".jsx", ".tsx", ".json", ".vue"]
    },
    define: {
        // __VUE_PROD_DEVTOOLS__: true,
        __VUE_PROD_DEVTOOLS__: mode !== "production",
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
    },
    customLogger: {
        info: message => {
            console.log(`[info] ${message}`)
        },
        warn: message => {
            console.log(`[warn] ${message}`)
        },
        error: message => {
            console.log(`[error] ${message}`)
        },
        debug: message => {
            console.log(`[debug] ${message}`)
        }
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
}))
