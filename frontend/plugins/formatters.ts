export default defineNuxtPlugin((NuxtApp) => {
  return {
    provide: {
      fullDate: function (dateStr: string | Date) {
        return fullDate(new Date(dateStr));
      },
      fullDotDate: function (dateStr: string | Date) {
        return fullDotDate(new Date(dateStr));
      },
      meetingDate: function (dateStr: string | Date) {
        return meetingDate(new Date(dateStr));
      },
      onlyDateDot: function (dateStr: string | Date) {
        return onlyDateDot(new Date(dateStr));
      },
      onlyDateSlash: function (dateStr: string | Date) {
        return onlyDateSlash(new Date(dateStr));
      },
      coolDate: function (dateStr: string | Date) {
        return coolDate(new Date(dateStr));
      },
      onlyDate: function (dateStr: string | Date) {
        return onlyDate(new Date(dateStr));
      },
      onlyTime: function (dateStr: string | Date) {
        return onlyTime(new Date(dateStr));
      },
      relativeDate: function (dateStr: string | Date) {
        return relativeDate(new Date(dateStr));
      },
      relativeFullDate: function (dateStr: string | Date) {
        return relativeFullDate(new Date(dateStr));
      },
      secondsToHms,
      formatDuration,
      humanSize,
      formatMoney,
      secondsToTime,
      convertToSeconds,
      convertTimeToSeconds,
      convertSecondsToTime,
      fileNameFormatRemove,
      toFixedWithoutZeros,
    },
  };
});

const GB = 1 << 30,
  MB = 1 << 20,
  KB = 1 << 10;

const humanSize = function (size: number): string {
  if (size > GB) {
    return (size / GB).toFixed(1) + "GB";
  } else if (size > MB) {
    return (size / MB).toFixed(1) + "MB";
  } else if (size > KB) {
    return (size / KB).toFixed(1) + "KB";
  } else {
    return size + "B";
  }
};
const toFixedWithoutZeros = (num: number, precision: number) =>
  num.toFixed(precision).replace(/\.0+$/, "");

function formatMoney(number: number): string {
  let formatted = "";
  if (isNaN(number)) {
    formatted = "";
  } else {
    formatted = Number(number).toLocaleString("mn-Mn", {
      maximumFractionDigits: 0,
      minimumFractionDigits: 0,
    });
  }
  return formatted;
}

function formatDuration(
  seconds: number,
  format: "object" | "mm:ss" | "hh:mm:ss" = "hh:mm:ss"
) {
  const h = Math.floor(seconds / (60 * 60))
    .toString()
    .padStart(2, "0");
  const m = Math.floor((seconds % (60 * 60)) / 60)
    .toString()
    .padStart(2, "0");
  const s = Math.floor(seconds % 60)
    .toString()
    .padStart(2, "0");

  if (format === "mm:ss") {
    return `${m}:${s}`;
  } else if (format === "hh:mm:ss") {
    return `${h}:${m}:${s}`;
  } else {
    return { h, m, s };
  }
}

const fullDate = function (date: Date) {
  return onlyDate(date).replaceAll("-", "/") + ` ${onlyTime(date)}`;
};

const fullDotDate = function (date: Date) {
  return onlyDate(date).replaceAll("-", ".") + ` ${onlyTime(date)}`;
};

const meetingDate = function (date: Date) {
  return onlyDate(date).replaceAll("-", "-") + ` ${onlyTime(date)}`;
};

const today = new Date();

const dateDiffInDays = (a: Date, b: Date) => {
  const MS_PER_DAY = 1000 * 60 * 60 * 24;
  const utc1 = Date.UTC(a.getFullYear(), a.getMonth(), a.getDate());
  const utc2 = Date.UTC(b.getFullYear(), b.getMonth(), b.getDate());

  return Math.floor((utc2 - utc1) / MS_PER_DAY);
};

const relativeDate = (date: Date) => {
  const diffDays = dateDiffInDays(today, date);
  let relativeDateStr = "";
  if (diffDays == 0) relativeDateStr += "Өнөөдөр";
  if (diffDays == 1) relativeDateStr += "Маргааш";
  if (diffDays == 2) relativeDateStr += "Нөгөөдөр";
  if (diffDays == -1) relativeDateStr += "Өчигдөр";
  if (diffDays == -2) relativeDateStr += "Уржигдар";
  if (diffDays > 2 || diffDays < -2)
    relativeDateStr += onlyDate(date).replaceAll("-", "/");

  return relativeDateStr;
};

const relativeFullDate = (date: Date) => {
  const diffDays = dateDiffInDays(today, date);
  let relativeDateStr = "";
  if (diffDays == 0) relativeDateStr += "Өнөөдөр";
  if (diffDays == 1) relativeDateStr += "Маргааш";
  if (diffDays == 2) relativeDateStr += "Нөгөөдөр";
  if (diffDays == -1) relativeDateStr += "Өчигдөр";
  if (diffDays == -2) relativeDateStr += "Уржигдар";
  if (diffDays > 2 || diffDays < -2)
    relativeDateStr += onlyDate(date).replaceAll("-", "/");

  return (relativeDateStr += ` ${onlyTime(date)}`);
};

const coolDate = function (date: Date) {
  const seconds = Math.abs(
    Math.floor((new Date().valueOf() - date.valueOf()) / 1000)
  );

  const isPast = new Date() > date;
  if (seconds < 86400) {
    let interval = Math.floor(seconds / 3600);
    if (interval >= 1) {
      if (interval == 1) {
        return "1 " + `цагийн ${isPast ? "өмнө" : "дараа"}`;
      }
      return interval + " " + `цагийн ${isPast ? "өмнө" : "дараа"}`;
    }
    interval = Math.floor(seconds / 60);
    if (interval >= 1) {
      if (interval == 1) {
        return "1 " + `минутын ${isPast ? "өмнө" : "дараа"}`;
      }
      return interval + " " + `минутын ${isPast ? "өмнө" : "дараа"}`;
    }
    return "Одоо";
  } else {
    return relativeDate(date);
  }
};

const onlyDate = function (date: Date) {
  return (
    date.getFullYear() +
    "-" +
    ("0" + (date.getMonth() + 1)).slice(-2) +
    "-" +
    ("0" + date.getDate()).slice(-2)
  );
};

const onlyDateDot = function (date: Date) {
  return (
    date.getFullYear() +
    "." +
    ("0" + (date.getMonth() + 1)).slice(-2) +
    "." +
    ("0" + date.getDate()).slice(-2)
  );
};

const onlyDateSlash = function (date: Date) {
  return (
    date.getFullYear() +
    "/" +
    ("0" + (date.getMonth() + 1)).slice(-2) +
    "/" +
    ("0" + date.getDate()).slice(-2)
  );
};

const onlyTime = function (date: Date) {
  let dateTime = new Date(date);
  let hour = dateTime.getHours();
  let minute = dateTime.getMinutes();
  let h = hour < 10 ? "0" + hour : hour.toString(),
    m = minute < 10 ? "0" + minute : minute.toString();
  return h + `:` + m;
};

function secondsToTime(secs: number) {
  const secondsInYear = 31536000;

  let years = Math.floor(secs / secondsInYear);
  var hours = Math.floor(secs / (60 * 60));

  var divisor_for_minutes = secs % (60 * 60);
  var minutes = Math.floor(divisor_for_minutes / 60);

  var divisor_for_seconds = divisor_for_minutes % 60;
  var seconds = Math.ceil(divisor_for_seconds);

  var obj = {
    y: years,
    h: hours,
    m: minutes,
    s: seconds,
    full: `${hours.toString().length == 1 ? "0" + hours : hours}:${
      minutes.toString().length == 1 ? "0" + minutes : minutes
    }:${seconds.toString().length == 1 ? "0" + seconds : seconds}`,
  };
  return obj;
}

function secondsToHms(seconds: number) {
  var h = Math.floor(seconds / (60 * 60));
  var m = Math.floor((seconds % (60 * 60)) / 60);
  var s = Math.floor(seconds % 60);
  var hDisplay = h > 0 ? h + (h == 1 ? " цаг " : " цаг ") : "";
  var mDisplay = m > 0 ? m + (m == 1 ? " мин " : " мин ") : "";
  var sDisplay = s > 0 ? s + (s == 1 ? " сек " : " сек ") : "";
  return hDisplay + mDisplay + sDisplay;
}

function convertToSeconds(hours: number, minutes: number, seconds: number) {
  return hours * 3600 + minutes * 60 + seconds;
}

function convertTimeToSeconds(value: number, type: string) {
  switch (type) {
    case "day":
      return Math.floor(value * 24 * 3600);
    case "h":
      return value * 60 * 60;
    case "m":
      return value * 60;
    case "s":
      return value;
    default:
      break;
  }
}

function convertSecondsToTime(value: number, type: string) {
  switch (type) {
    case "day":
      return Math.floor(value / 86400);
    case "h":
      return value / 3600;
    case "m":
      return value / 60;
    case "s":
      return value;
    default:
      break;
  }
}

function fileNameFormatRemove(fileName: string): string {
  const lastDotIndex = fileName.lastIndexOf(".");

  if (lastDotIndex === -1) {
    return fileName;
  }

  return fileName.slice(0, 25);
}
