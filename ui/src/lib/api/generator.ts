import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { convertGenerators, Generator, GeneratorSchema } from "../types/generator"
import { convertTracks } from "../types/track"
import { apiDelete, apiGet, apiPost, apiPut } from "./query"
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

export const useGeneratorRefresh = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (generator: Pick<Generator, "id">) => apiPost(`${ENDPOINT}/refresh/${generator.id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["task"] })
      queryClient.invalidateQueries({ queryKey: ["generator"] })
    },
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

export const useGeneratorDelete = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (args: { generator: Pick<Generator, "id">, deletePlaylist?: boolean }) => {
      const queryParams = new URLSearchParams()

      if (args.deletePlaylist !== undefined) {
        queryParams.append("delete_playlist", String(args.deletePlaylist))
      }

      return apiDelete(`${ENDPOINT}/${args.generator.id}?${queryParams.toString()}`)
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["generator"] }),
  })
}
