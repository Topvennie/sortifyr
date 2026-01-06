import z from "zod";
import { JSONBody } from "./general";
import { API } from "./api";
import { convertTracks, Track } from "./track";

export enum GeneratorPreset {
  Top = "top",
  OldTop = "old_top"
}
export const generatorPresetString: Record<GeneratorPreset, string> = {
  [GeneratorPreset.Top]: "Top",
  [GeneratorPreset.OldTop]: "Old Top",
}

export interface GeneratorWindow {
  start: Date;
  end: Date;
  minPlays: number;
  burstIntervalDays: number;
}

export const convertGeneratorWindow = (g: API.GeneratorWindow): GeneratorWindow => {
  return {
    start: new Date(g.start),
    end: new Date(g.end),
    minPlays: g.min_plays,
    burstIntervalDays: g.burst_interval_days,
  }
}

export interface GeneratorParams {
  trackAmount: number;
  excludedPlaylistIds: number[];
  excludedTrackIds: number[];
  preset: GeneratorPreset;
  paramsTop?: {
    window: GeneratorWindow;
  };
  paramsOldTop?: {
    peakWindow: GeneratorWindow;
    recentWindow: GeneratorWindow;
  };
}

export const convertGeneratorParams = (g: Pick<API.Generator, "params">): GeneratorParams => {
  return {
    trackAmount: g.params.track_amount,
    excludedPlaylistIds: g.params.excluded_playlist_ids ?? [],
    excludedTrackIds: g.params.excluded_track_ids ?? [],
    preset: g.params.preset as GeneratorPreset,
    paramsTop: g.params.params_top ? {
      window: convertGeneratorWindow(g.params.params_top.window)
    } : undefined,
    paramsOldTop: g.params.params_old_top ? {
      peakWindow: convertGeneratorWindow(g.params.params_old_top.peak_window),
      recentWindow: convertGeneratorWindow(g.params.params_old_top.recent_window),
    } : undefined,
  }
}

export interface Generator {
  id: number;
  name: string;
  description?: string;
  playlistId?: number;
  intervalDays: number;
  spotifyOutdated: boolean;
  params: GeneratorParams;
  tracks: Track[];
  lastUpdate?: Date;
}

export const convertGenerator = (g: API.Generator): Generator => {
  return {
    id: g.id,
    name: g.name,
    description: g.description,
    playlistId: g.playlist_id,
    intervalDays: g.interval_days,
    spotifyOutdated: g.spotify_outdated,
    params: convertGeneratorParams(g),
    tracks: convertTracks(g.tracks),
    lastUpdate: g.last_update ? new Date(g.last_update) : undefined,
  }
}

export const convertGenerators = (g: API.Generator[]): Generator[] => {
  return g.map(convertGenerator)
}

export const convertGeneratorSchema = (g: Generator): GeneratorSchema => {
  return {
    id: g.id,
    name: g.name,
    description: g.description,
    createPlaylist: (g.playlistId ?? 0) !== 0,
    intervalDays: g.intervalDays,
    params: { ...g.params },
  }
}

export const generatorWindowSchema = z.object({
  start: z.date(),
  end: z.date(),
  minPlays: z.number().positive(),
  burstIntervalDays: z.number().positive(),
})
export type GeneratorWindowSchema = z.infer<typeof generatorWindowSchema> & JSONBody;

export const generatorParamsSchema = z.object({
  trackAmount: z.number().positive(),
  excludedPlaylistIds: z.array(z.number().positive()),
  excludedTrackIds: z.array(z.number().positive()),
  preset: z.enum(GeneratorPreset),
  paramsTop: z.object({
    window: generatorWindowSchema
  }).partial().optional(),
  paramsOldTop: z.object({
    peakWindow: generatorWindowSchema,
    recentWindow: generatorWindowSchema
  }).partial().optional(),
})
export type GeneratorParamsSchema = z.infer<typeof generatorParamsSchema> & JSONBody;

export const generatorSchema = z.object({
  id: z.number().positive().optional(),
  name: z.string().nonempty(),
  description: z.string().optional(),
  createPlaylist: z.boolean(),
  intervalDays: z.number().nonnegative(),
  params: generatorParamsSchema.partial().optional(),
})
export type GeneratorSchema = z.infer<typeof generatorSchema> & JSONBody;
