import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { convertGenerators, GeneratorSchema } from "../types/generator"
import { convertTracks } from "../types/track"
import { apiGet, apiPost, apiPut } from "./query"
import { STALE_TIME } from "../types/staletime"

const ENDPOINT = "generator"

export const useGeneratorGetAll = () => {
  return useQuery({
    queryKey: ["generator"],
    queryFn: async () => (await apiGet(ENDPOINT, convertGenerators)).data,
    staleTime: STALE_TIME.MIN_30,
    throwOnError: true,
  })
}

export const useGeneratorPreview = () => {
  return useMutation({
    mutationFn: async (generator: GeneratorSchema) => (await apiPost(`${ENDPOINT}/preview`, generator.params, convertTracks)).data,
    throwOnError: true,
  })
}

export const useGeneratorCreate = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (generator: GeneratorSchema) => apiPut(ENDPOINT, generator),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["generator"] }),
  })
}

export const useGeneratorEdit = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (generator: GeneratorSchema) => apiPost(`${ENDPOINT}/${generator.id}`, generator),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["generator"] }),
  })
}
