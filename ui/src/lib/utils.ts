import { type ClassValue, clsx } from "clsx";
import { format } from "date-fns";
import { twMerge } from "tailwind-merge";
import { v4 as uuid } from 'uuid';
import { isResponseNot200Error } from "./api/query";

export function camelToSnake(obj: unknown): unknown {
  if (obj === null || obj === undefined) {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(camelToSnake);
  }

  if (obj instanceof Date) {
    return obj.toISOString()
  }

  if (typeof obj === "object") {
    return Object.fromEntries(
      Object.entries(obj).map(([key, value]) => [
        stringCamelToSnake(key),
        camelToSnake(value),
      ]),
    );
  }

  return obj;
}

function stringCamelToSnake(str: string) {
  return str.replace(/[A-Z]+/g, l => `_${l.toLowerCase()}`);
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getUuid() {
  return uuid();
}

export function formatDate(date: Date) {
  return format(date, "eee dd MMMM, HH:mm");
}

export function debounce<T extends (...args: unknown[]) => void>(fn: T, delay: number) {
  let timer: ReturnType<typeof setTimeout>;

  return (...args: Parameters<T>): void => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  };
}

export function capitalize(text: string) {
  if (text.length <= 1) return text.toUpperCase()

  return text[0].toUpperCase() + text.slice(1)
}

export function getBuildTime() {
  const buildTime = import.meta.env.VITE_BUILD_TIME as string | "";

  return buildTime ? formatDate(new Date(buildTime)) : "Development"
}

export const scrollTo = (id: string) => {
  const element = document.getElementById(id)
  if (!element) return

  element.scrollIntoView({ behavior: "smooth" })
}

export async function getErrorMessage(err: Error) {
  if (isResponseNot200Error(err)) return await err.response.text()

  return err.message
}

export function daysAgo(days: number) {
  return new Date(Date.now() - days * 24 * 60 * 60 * 1000);
}

export function getValueByPath<T>(obj: unknown, path: string): T | undefined {
  return path
    .split(".")
    .reduce<unknown>((acc, key) => {
      if (acc == null || typeof acc !== "object") {
        return undefined;
      }
      return (acc as Record<string, unknown>)[key];
    }, obj) as T | undefined;
}
