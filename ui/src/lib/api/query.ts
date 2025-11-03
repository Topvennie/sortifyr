

import { QueryClient } from "@tanstack/react-query";
import { camelToSnake } from "../utils";
import { JSONBody } from "../types/general";
import { CONTENT_TYPE } from "../types/contentType";

export const queryClient = new QueryClient();

export const NO_DATA: JSONBody = {}
export const NO_CONVERTER = undefined
export const NO_FILES = undefined

export type FileData = {
  file: File
  field: string
}

const URLS: Record<string, string> = {
  PUBLIC: "/api",
};

export async function apiGet<T, U = unknown>(
  endpoint: string,
  convertData?: (data: U) => T,
) {
  return _fetch<T, U>(
    `${URLS.PUBLIC}/${endpoint}`,
    {},
    convertData,
  );
}

export async function apiPost<T, U = unknown>(
  endpoint: string,
  data: JSONBody = {},
  convertData?: (data: U) => T,
  files?: FileData[],
) {
  const { headers, body } = _buildFormData(data, files)

  return _fetch<T, U>(
    `${URLS.PUBLIC}/${endpoint}`,
    {
      method: "POST",
      headers,
      body,
    },
    convertData,
  );
}

export async function apiPut<T, U = unknown>(
  endpoint: string,
  data: JSONBody = {},
  convertData?: (data: U) => T,
  files?: FileData[],
) {
  const { headers, body } = _buildFormData(data, files)

  return _fetch<T, U>(
    `${URLS.PUBLIC}/${endpoint}`,
    {
      method: "PUT",
      headers,
      body
    },
    convertData,
  );
}

export async function apiPatch<T, U = unknown>(
  endpoint: string,
  data: JSONBody = {},
  convertData?: (data: U) => T,
  files?: FileData[],
) {
  const { headers, body } = _buildFormData(data, files)

  return _fetch<T, U>(
    `${URLS.PUBLIC}/${endpoint}`,
    {
      method: "PATCH",
      headers,
      body
    },
    convertData,
  );
}

export async function apiDelete<T, U = unknown>(
  endpoint: string,
  convertData?: (data: U) => T,
) {
  return _fetch<T, U>(
    `${URLS.PUBLIC}/${endpoint}`,
    {
      method: "DELETE",
    },
    convertData,
  );
}

interface ResponseNot200Error extends Error {
  response: Response;
}

export function isResponseNot200Error(error: unknown): error is ResponseNot200Error {
  return (error as ResponseNot200Error).response !== undefined;
}

function _buildFormData(data: JSONBody, files?: FileData[]): RequestInit {
  if (!files?.length) {
    return {
      body: JSON.stringify(camelToSnake(data)),
      headers: {
        "Content-Type": CONTENT_TYPE.JSON
      }
    }
  }

  const formData = new FormData();

  Object.entries(camelToSnake(data) as Record<string, unknown>).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      if (typeof value === "string") {
        formData.append(key, value)
      } else {
        formData.append(key, JSON.stringify(value));
      }
    }
  });

  files.forEach(f => formData.append(f.field, f.file))

  return {
    body: formData,
  };
}

async function _fetch<T, U>(url: string, options: RequestInit = {}, convertData?: (data: U) => T): Promise<{ data: T, response: Response }> {
  return fetch(
    url,
    { credentials: "include", ...options },
  ).then(async (response) => {
    if (!response.ok) {
      const error = new Error(`Fetch failed with status: ${response.status}`) as ResponseNot200Error;
      error.response = response;
      throw error;
    }

    const contentType = response.headers.get("content-type");

    if (contentType?.includes(CONTENT_TYPE.JSON))
      return { data: await response.json() as Promise<unknown>, response };
    else if ([CONTENT_TYPE.PDF, CONTENT_TYPE.PNG, CONTENT_TYPE.WEBP].some(t => contentType?.includes(t)))
      return { data: await response.blob(), response };
    else
      return { data: await response.text(), response };
  }).then(({ data, response }: { data: unknown, response: Response }) => ({
    data: convertData ? convertData(data as U) : (data as T), response
  }));
}
