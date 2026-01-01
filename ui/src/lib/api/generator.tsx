import { useQuery } from "@tanstack/react-query"
import { GeneratorSchema } from "../types/generator"
import { apiPost } from "./query"
import { convertTracks } from "../types/track"

const ENDPOINT = "generator"

export const useGeneratorGenerate = (generator: GeneratorSchema) => {
  return useQuery({
    queryKey: ["generator", "generate"],
    queryFn: async () => (await apiPost(`${ENDPOINT}/generate`, generator, convertTracks)).data,
    staleTime: Infinity,
    throwOnError: true,
  })
}
