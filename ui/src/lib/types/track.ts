import { API } from "./api";
import { convertPlaylist, Playlist } from "./playlist";

export interface Track {
  id: number;
  spotifyId: string;
  name: string;
}

export const convertTrack = (t: API.Track): Track => {
  return {
    id: t.id,
    spotifyId: t.spotify_id,
    name: t.name,
  }
}

export interface TrackHistory extends Track {
  historyId: number;
  playedAt: Date;
  playCount?: number;
}

export const convertTrackHistory = (h: API.TrackHistory): TrackHistory => {
  return {
    ...convertTrack(h),
    historyId: h.history_id,
    playedAt: new Date(h.played_at),
    playCount: h.play_count,
  }
}

export const convertTrackHistories = (h: API.TrackHistory[]): TrackHistory[] => {
  return h.map(convertTrackHistory)
}

export interface TrackAdded extends Track {
  playlist: Playlist;
  createdAt: Date;
}

export const convertTrackAdded = (t: API.TrackAdded): TrackAdded => {
  return {
    ...convertTrack(t),
    playlist: convertPlaylist(t.playlist),
    createdAt: new Date(t.created_at),
  }
}

export const convertTracksAdded = (t: API.TrackAdded[]): TrackAdded[] => {
  return t.map(convertTrackAdded)
}

export interface TrackDeleted extends Track {
  playlist: Playlist;
  deletedAt: Date;
}

export const convertTrackDeleted = (t: API.TrackDeleted): TrackDeleted => {
  return {
    ...convertTrack(t),
    playlist: convertPlaylist(t.playlist),
    deletedAt: new Date(t.deleted_at),
  }
}

export const convertTracksDeleted = (t: API.TrackDeleted[]): TrackDeleted[] => {
  return t.map(convertTrackDeleted)
}

export interface TrackDuplicate extends Track {
  amount: number;
}

export const convertTrackDuplicate = (t: API.TrackDuplicate): TrackDuplicate => {
  return {
    ...convertTrack(t),
    amount: t.amount,
  }
}

export interface TrackFilter {
  playlistId?: string;
}

export interface TrackHistoryFilter {
  skipped?: boolean;
  start?: Date;
  end?: Date;
}
