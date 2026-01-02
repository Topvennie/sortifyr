import { useMutation } from "@tanstack/react-query"
import { GeneratorSchema } from "../types/generator"
import { convertTracks } from "../types/track"
import { apiPost } from "./query"

const ENDPOINT = "generator"

export const useGeneratorPreview = () => {
  return useMutation({
    mutationFn: async (generator: GeneratorSchema) => (await apiPost(`${ENDPOINT}/preview`, generator.params, convertTracks)).data,
    throwOnError: true,
  })
}
