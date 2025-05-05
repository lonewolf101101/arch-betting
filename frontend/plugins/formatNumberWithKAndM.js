export default defineNuxtPlugin((nuxtApp) => {
  var convert = function (number) {
    if (isNaN(number)) {
      return "Invalid input";
    }

    if (number < 1000) {
      return number.toString();
    } else if (number < 1000000) {
      return (number / 1000).toFixed(1) + "K";
    } else {
      return (number / 1000000).toFixed(1) + "M";
    }
  };

  return {
    provide: {
      formatNumberWithKAndM: (number) => {
        return convert(number);
      },
    },
  };
});
