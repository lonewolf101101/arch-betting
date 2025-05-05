export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook("app:mounted", () => {
    trackPageVisit();
  });

  const router = useRouter();
  router.afterEach((to) => {
    trackPageVisit(to.path);
  });

  // Function to send visit data to backend
  const trackPageVisit = async (path = window.location.pathname) => {
    try {
      await fetch("localhost:8080/pub/visit", {
        keepalive: true,
      });
    } catch (error) {
      console.error("Failed to track page visit:", error);
    }
  };
});
