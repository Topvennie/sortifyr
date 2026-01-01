import z from "zod";
import { JSONBody } from "./general";

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

export interface GeneratorParams {
  trackAmount: number;
  minPlayCount: number;
}

export interface Generator {
  preset: GeneratorPreset;
  params: GeneratorParams;
}

export const generatorParamsSchema = z.object({
  trackAmount: z.number().positive(),
  minPlayCount: z.number().positive(),
})
export type GeneratorParamsSchema = z.infer<typeof generatorParamsSchema> & JSONBody;

export const generatorSchema = z.object({
  preset: z.enum(GeneratorPreset),
  params: generatorParamsSchema.partial(),
})
export type GeneratorSchema = z.infer<typeof generatorSchema> & JSONBody;
