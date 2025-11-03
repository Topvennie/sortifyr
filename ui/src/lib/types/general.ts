export interface JSON {
  [key: string]: string | number | boolean | JSON | JSON[] | unknown;
}

export type JSONBody = JSON;
