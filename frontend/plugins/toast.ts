import * as vt from "vue-toastification";
import "vue-toastification/dist/index.css";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(vt.default, {
    hideProgressBar: true,
    position: vt.POSITION.TOP_CENTER,
    timeout: 2000,
    closeButton: false,
  });
  return {
    provide: {
      toast: vt.useToast(),
    },
  };
});
