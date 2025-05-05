// https://nuxt.com/docs/api/configuration/nuxt-config
import vuetify, { transformAssetUrls } from "vite-plugin-vuetify";

const isDev = process.env.NODE_ENV !== "production";

const domain = process.env.DOMAIN;

console.log("Server type:", process.env.SERVER_TYPE);

const nitro =
  process.env.SERVER_TYPE != "local"
    ? {
        routeRules: {
          "pub/**": { proxy: "http://localhost:3000/api/**" },
        },
      }
    : {
        routeRules: {
          "/api/**": { proxy: "http://localhost:8080/api/**" },
          "/pub/**": { proxy: "http://localhost:8080/pub/**" },
        },
        devProxy: {
          "/api": {
            target: "http://localhost:8080/api/",
            changeOrigin: true,
            ws: true,
          },
          "/pub": {
            target: "http://localhost:8080/pub/",
            changeOrigin: true,
            ws: true,
          },
        },
      };

export default defineNuxtConfig({
  devtools: { enabled: true },

  build: {
    transpile: ["vuetify"],
  },

  plugins: [
    "~/plugins/toast.ts",
    "~/plugins/formatNumberWithKAndM.js",
    "~/plugins/mitt.ts",
    "~/plugins/visit-tracker.js",
  ],

  modules: [
    (_options, nuxt) => {
      nuxt.hooks.hook("vite:extendConfig", (config) => {
        // @ts-expect-error
        config.plugins.push(
          vuetify({
            autoImport: true,
            styles: {
              configFile: "./assets/settings.scss",
            },
          })
        );
      });
    },
  ],

  ssr: false,

  sourcemap: {
    server: false,
    client: false,
  },

  runtimeConfig: {
    public: {
      dev: isDev,
      domain: domain,
    },
  },

  vite: {
    vue: {
      template: {
        transformAssetUrls,
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: ["legacy-js-api"],
        },
        sass: {
          silenceDeprecations: ["legacy-js-api"],
        },
      },
    },
  },

  nitro,

  app: {
    head: {
      title: "Betting",
      titleTemplate: "%s",
    },
  },

  compatibilityDate: "2025-05-05",
});
