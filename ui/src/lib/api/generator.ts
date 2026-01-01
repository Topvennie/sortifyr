import { useQuery } from "@tanstack/react-query"
import { GeneratorSchema } from "../types/generator"
import { apiPost } from "./query"
import { convertTracks } from "../types/track"

const ENDPOINT = "generator"

export const useGeneratorPreview = (generator: GeneratorSchema) => {
  return useQuery({
    queryKey: ["generator", "preview"],
    queryFn: async () => (await apiPost(`${ENDPOINT}/preview`, generator.params, convertTracks)).data,
    staleTime: Infinity,
    throwOnError: true,
  })
}
