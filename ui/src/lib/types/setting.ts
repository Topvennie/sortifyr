import { API } from "./api";

export interface Setting {
  id: number;
  lastUpdate?: Date;
}

export const convertSetting = (s: API.Setting): Setting => {
  return {
    id: s.id,
    lastUpdate: s.last_update ? new Date(s.last_update) : undefined,
  }
}
