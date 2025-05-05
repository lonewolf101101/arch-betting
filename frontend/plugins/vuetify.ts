import { mn } from "date-fns/locale/mn";
import "vuetify/styles";
import { aliases, fr } from "../iconsets/feather-iconset";
import { createVuetify } from "vuetify";

export default defineNuxtPlugin((app) => {
  const vuetify = createVuetify({
    defaults: {
      VBtn: {
        style: "text-transform: none; letter-spacing: 0.005em;",
      },
      VAppBar: {
        height: 80,
      },
      global: {
        style: {
          fontFamily: "Rubik, sans-serif",
        },
      },
    },
    theme: {
      themes: {
        light: {
          dark: false,
          colors: {
            success: "#4BE262",
            primary: "#EDEEEF",
            // primary: "#2060eb",
            darkPrimary: "#282929",
            secondary: "#003BFD",
            grey: "#656567",
            gray100: "#474748",
            black: "#6B6C6E",
            brown: "#606162",
            error: "#ec5551",
            warning: "#F1AA09",
            warningText: "#863A00",
            textDark: "#0522A0",
            neutral: "#fbfcfe",
            darkBlue: "#0828D8",
          },
        },
      },
    },
    date: {
      locale: {
        mn,
      },
    },
    icons: {
      defaultSet: "fr",
      aliases,
      sets: {
        fr,
      },
    },
  });
  app.vueApp.use(vuetify);
});
