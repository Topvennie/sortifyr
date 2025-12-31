import { useInfiniteQuery } from "@tanstack/react-query"
import { STALE_TIME } from "../types/staletime"
import { convertTrackHistories, convertTracksAdded, convertTracksDeleted, TrackFilter, TrackHistoryFilter } from "../types/track"
import { apiGet } from "./query"

const ENDPOINT = "track"
const PAGE_LIMIT = 100

export const useTrackGetHistory = (filter?: TrackHistoryFilter) => {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["track", "history", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      })

      if (filter?.skipped !== undefined) {
        queryParams.append("skipped", String(filter?.skipped))
      }
      if (filter?.start !== undefined) {
        queryParams.append("start", filter.start.toISOString())
      }
      if (filter?.end !== undefined) {
        queryParams.append("end", filter.end.toISOString())
      }

      const url = `${ENDPOINT}/history?${queryParams.toString()}`
      return (await apiGet(url, convertTrackHistories)).data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1
    },
    staleTime: STALE_TIME.MIN_5,
    throwOnError: true,
  })

  const history = data?.pages.flat() ?? []

  return {
    history,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  }
}

export const useTrackGetAdded = (filter?: TrackFilter) => {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["track", "added", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      })

      if (filter?.playlistId !== undefined) {
        queryParams.append("playlist_id", filter.playlistId)
      }

      const url = `${ENDPOINT}/added?${queryParams.toString()}`
      return (await apiGet(url, convertTracksAdded)).data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1
    },
    staleTime: STALE_TIME.MIN_30,
    throwOnError: true,
  })

  const tracks = data?.pages.flat() ?? []

  return {
    tracks,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  }
}

export const useTrackGetDeleted = (filter?: TrackFilter) => {
  const { data, isLoading, fetchNextPage, isFetchingNextPage, hasNextPage, error, refetch, isFetching } = useInfiniteQuery({
    queryKey: ["track", "deleted", JSON.stringify(filter)],
    queryFn: async ({ pageParam = 1 }) => {
      const queryParams = new URLSearchParams({
        page: pageParam.toString(),
        limit: PAGE_LIMIT.toString(),
      })

      if (filter?.playlistId !== undefined) {
        queryParams.append("playlist_id", filter.playlistId)
      }

      const url = `${ENDPOINT}/deleted?${queryParams.toString()}`
      return (await apiGet(url, convertTracksDeleted)).data
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      return lastPage.length < PAGE_LIMIT ? undefined : allPages.length + 1
    },
    staleTime: STALE_TIME.MIN_30,
    throwOnError: true,
  })

  const tracks = data?.pages.flat() ?? []

  return {
    tracks,
    isLoading,
    fetchNextPage,
    isFetchingNextPage,
    hasNextPage,
    error,
    refetch,
    isFetching,
  }
}
