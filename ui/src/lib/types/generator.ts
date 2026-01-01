import z from "zod";
import { JSONBody } from "./general";
import { API } from "./api";

export enum GeneratorPreset {
  Custom = "custom",
  Forgotten = "forgotten",
  Top = "top",
  OldTop = "old_top"
}
export const generatorPresetString: Record<GeneratorPreset, string> = {
  [GeneratorPreset.Custom]: "Custom",
  [GeneratorPreset.Forgotten]: "Forgotten",
  [GeneratorPreset.Top]: "Top",
  [GeneratorPreset.OldTop]: "Old Top",
}

export interface GeneratorWindow {
  start: Date;
  end: Date;
  minPlays: number;
  burstIntervalS: number;
}

export const convertGeneratorWindow = (g: API.GeneratorWindow): GeneratorWindow => {
  return {
    start: new Date(g.start),
    end: new Date(g.end),
    minPlays: g.min_plays,
    burstIntervalS: g.burst_interval_s,
  }
}

export interface Generator {
  id: number;
  name: string;
  description?: string;
  params: {
    trackAmount: number;
    excludedPlaylistIds: number[];
    excludedTrackIds: number[];
    preset: GeneratorPreset;
    paramsCustom?: {};
    paramsForgotten?: {};
    paramsTop?: {
      window: GeneratorWindow;
    };
    paramsOldTop?: {
      peakWindow: GeneratorWindow;
      recentWindow: GeneratorWindow;
    };
  };
}

export const convertGenerator = (g: API.Generator): Generator => {
  return {
    id: g.id,
    name: g.name,
    description: g.description,
    params: {
      trackAmount: g.params.track_amount,
      excludedPlaylistIds: g.params.excluded_playlist_ids ?? [],
      excludedTrackIds: g.params.excluded_track_ids ?? [],
      preset: g.params.preset as GeneratorPreset,
      paramsCustom: g.params.params_custom,
      paramsForgotten: g.params.params_forgotten,
      paramsTop: g.params.params_top ? {
        window: convertGeneratorWindow(g.params.params_top.window)
      } : undefined,
      paramsOldTop: g.params.params_old_top ? {
        peakWindow: convertGeneratorWindow(g.params.params_old_top.peak_window),
        recentWindow: convertGeneratorWindow(g.params.params_old_top.current_window),
      } : undefined,
    }
  }
}

export const convertGeneratorSchema = (g: Generator): GeneratorSchema => {
  return {
    ...g
  }
}

export const generatorWindowSchema = z.object({
  start: z.date(),
  end: z.date(),
  minPlays: z.number().positive(),
  burstIntervalS: z.number().positive(),
})
export type GeneratorWindowSchema = z.infer<typeof generatorWindowSchema> & JSONBody;

export const generatorSchema = z.object({
  id: z.number().positive().optional(),
  name: z.string().nonempty(),
  description: z.string().optional(),
  params: z.object({
    trackAmount: z.number().positive(),
    excludedPlaylistIds: z.array(z.number().positive()).optional(),
    excludedTrackIds: z.array(z.number().positive()).optional(),
    preset: z.enum(GeneratorPreset),
    paramsCustom: z.object({}).partial().optional(),
    paramsForgotten: z.object({}).partial().optional(),
    paramsTop: z.object({
      window: generatorWindowSchema
    }).partial().optional(),
    paramsOldTop: z.object({
      peakWindow: generatorWindowSchema,
      recentWindow: generatorWindowSchema
    }).partial().optional(),
  }).partial().optional(),
})
export type GeneratorSchema = z.infer<typeof generatorSchema> & JSONBody;
