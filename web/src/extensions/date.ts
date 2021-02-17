import {
  addDays,
  addMinutes,
  addMonths,
  addSeconds,
  addWeeks,
  addYears,
  differenceInCalendarDays,
  endOfISOWeek,
  endOfMonth,
  format as dateFnsFormat,
  isSameDay,
  isTomorrow,
  isYesterday,
  startOfISOWeek,
  startOfMonth,
} from "date-fns";
import { fr } from "date-fns/locale";

/**
 * DateTime is just a TS type around the ISO strings returned by the backend
 * It avoids transforming them manually by calling a constructor each time
 * Or create a complex system of codegen + Apollo link to transform them on the fly
 */

declare global {
  interface DateTime {
    getDate(this: DateTime): Date;
    getTime(this: DateTime): number;

    millisecondsFromNow(this: DateTime): number;
    secondsFromNow(this: DateTime): number;

    plusSeconds(this: DateTime, seconds: number): DateTime;
    minusSeconds(this: DateTime, seconds: number): DateTime;
    plusMinutes(this: DateTime, minutes: number): DateTime;
    minusMinutes(this: DateTime, minutes: number): DateTime;
    plusDays(this: DateTime, days: number): DateTime;
    minusDays(this: DateTime, days: number): DateTime;
    plusWeeks(this: DateTime, weeks: number): DateTime;
    minusWeeks(this: DateTime, weeks: number): DateTime;
    plusMonths(this: DateTime, months: number): DateTime;
    minusMonths(this: DateTime, months: number): DateTime;
    plusYears(this: DateTime, years: number): DateTime;
    minusYears(this: DateTime, years: number): DateTime;

    startOfWeek(this: DateTime): DateTime;
    endOfWeek(this: DateTime): DateTime;
    startOfMonth(this: DateTime): DateTime;
    endOfMonth(this: DateTime): DateTime;

    isSameDay(this: DateTime, date: DateTime): boolean;
    isToday(this: DateTime): boolean;
    isBefore(this: DateTime, date: DateTime): boolean;
    isAfter(this: DateTime, date: DateTime): boolean;
    isPast(this: DateTime): boolean;
    isFuture(this: DateTime): boolean;
    isBeforeToday(this: DateTime): boolean;
    isAfterToday(this: DateTime): boolean;

    format(
      this: DateTime,
      format:
        | "date" // 25/01/2021
        | "monthDay" // 25 (for calendars)
        | "time" // 22:46
        | "monthAndYear" // janv. 2021
        | {
            relative:
              | "timeOrDate" // il y a 10m, il y a 3h, {week if >1d}
              | "week" // aujourd'hui, hier, mercredi, {date}
              | "weekWithTime"; // {week} à {time}
            // Adds "le" before non relative date
            withDatePrefix?: boolean;
          }
        // For other exceptional cases (calendars, pickers, files)
        // See https://date-fns.org/v2.16.1/docs/format
        | { exception: string },
    ): string;
  }

  interface Date {
    toDateTime(): DateTime;
  }
}

Object.defineProperty(Date.prototype, "toDateTime", {
  value() {
    return this.toISOString();
  },
});

const now = () => new Date().toDateTime();

const defineDateTimeProperty = <K extends keyof DateTime>(
  key: K,
  value: DateTime[K],
) => {
  Object.defineProperty(String.prototype, key, { value });
};

defineDateTimeProperty("getDate", function () {
  return new Date((this as unknown) as string);
});
defineDateTimeProperty("getTime", function () {
  return this.getDate().getTime();
});

defineDateTimeProperty("millisecondsFromNow", function () {
  return Math.abs(new Date().getTime() - this.getTime());
});
defineDateTimeProperty("secondsFromNow", function () {
  return Math.floor(this.millisecondsFromNow() / 1000);
});

defineDateTimeProperty("plusSeconds", function (seconds) {
  return addSeconds(this.getDate(), seconds).toDateTime();
});
defineDateTimeProperty("minusSeconds", function (seconds) {
  return this.plusSeconds(-seconds);
});
defineDateTimeProperty("plusMinutes", function (minutes) {
  return addMinutes(this.getDate(), minutes).toDateTime();
});
defineDateTimeProperty("minusMinutes", function (minutes) {
  return this.plusMinutes(-minutes);
});
defineDateTimeProperty("plusDays", function (days) {
  return addDays(this.getDate(), days).toDateTime();
});
defineDateTimeProperty("minusDays", function (days) {
  return this.plusDays(-days);
});
defineDateTimeProperty("plusWeeks", function (weeks) {
  return addWeeks(this.getDate(), weeks).toDateTime();
});
defineDateTimeProperty("minusWeeks", function (weeks) {
  return this.plusWeeks(-weeks);
});
defineDateTimeProperty("plusMonths", function (months) {
  return addMonths(this.getDate(), months).toDateTime();
});
defineDateTimeProperty("minusMonths", function (months) {
  return this.plusMonths(-months);
});
defineDateTimeProperty("plusYears", function (years) {
  return addYears(this.getDate(), years).toDateTime();
});
defineDateTimeProperty("minusYears", function (years) {
  return this.plusYears(-years);
});

defineDateTimeProperty("startOfWeek", function () {
  // Until we target US/Canada, startOfWeek = startOfISOWeek
  return startOfISOWeek(this.getDate()).toDateTime();
});
defineDateTimeProperty("endOfWeek", function () {
  return endOfISOWeek(this.getDate()).toDateTime();
});
defineDateTimeProperty("startOfMonth", function () {
  return startOfMonth(this.getDate()).toDateTime();
});
defineDateTimeProperty("endOfMonth", function () {
  return endOfMonth(this.getDate()).toDateTime();
});

defineDateTimeProperty("isSameDay", function (date) {
  return isSameDay(this.getDate(), date.getDate());
});
defineDateTimeProperty("isToday", function () {
  return this.isSameDay(now());
});
defineDateTimeProperty("isBefore", function (date) {
  return this.getTime() < date.getTime();
});
defineDateTimeProperty("isAfter", function (date) {
  return this.getTime() > date.getTime();
});
defineDateTimeProperty("isPast", function () {
  return this.isBefore(now());
});
defineDateTimeProperty("isFuture", function () {
  return this.isAfter(now());
});
defineDateTimeProperty("isBeforeToday", function () {
  return this.isPast() && !this.isToday();
});
defineDateTimeProperty("isAfterToday", function () {
  return this.isFuture() && !this.isToday();
});

defineDateTimeProperty("format", function (format) {
  const libFormat = (stringFormat: string) =>
    dateFnsFormat(this.getDate(), stringFormat, { locale: fr });

  if (typeof format === "string") {
    return libFormat(
      {
        monthAndYear: "MMM y",
        date: "dd/MM/yyyy",
        monthDay: "d",
        time: "HH:mm",
      }[format],
    );
  }

  if ("exception" in format) return libFormat(format.exception);

  // Relative
  if (format.relative === "timeOrDate") {
    const numSecs = this.secondsFromNow();
    if (numSecs < 60) return "< 1mn";
    if (numSecs < 60 * 60) return `il y a ${Math.floor(numSecs / 60)}mn`;
    if (numSecs < 24 * 60 * 60) {
      return `il y a ${Math.floor(numSecs / 60 / 60)}h`;
    }
  }

  const relativeDate = (() => {
    if (this.isToday()) return "aujourd'hui";
    if (isYesterday(this.getDate())) return "hier";
    if (isTomorrow(this.getDate())) return "demain";
    if (Math.abs(differenceInCalendarDays(new Date(), this.getDate())) < 7) {
      return this.format({ exception: "EEEE" });
    }
    return `${format.withDatePrefix ? "le " : ""}${this.format("date")}`;
  })();

  return {
    weekWithTime: `${relativeDate} à ${this.format("time")}`,
    timeOrDate: relativeDate,
    week: relativeDate,
  }[format.relative];
});
