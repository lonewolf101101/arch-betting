import mitt from "mitt";

export default defineNuxtPlugin((nuxtApp) => {
  const emitter = mitt();

  nuxtApp.provide("emitEvent", emitter.emit); // Emit events
  nuxtApp.provide("listen", emitter.on); // Listen for events
  nuxtApp.provide("clear", emitter.off); // Clear listeners
});
