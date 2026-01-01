export namespace API {
  interface Base extends JSON {
    id: number;
  }

  export interface User extends Base {
    uid: string;
    name: string;
    display_name: string;
    email: string;
  }

  export interface Playlist extends Base {
    spotify_id: string;
    owner?: User;
    name: string;
    description?: string;
    public?: boolean;
    track_amount: number;
    collaborative?: boolean;
    has_cover: boolean;
  }

  export interface PlaylistDuplicate extends Playlist {
    duplicates: TrackDuplicate[],
  }

  export interface PlaylistUnplayable extends Playlist {
    unplayables: Track[],
  }

  export interface Directory extends Base {
    name: string;
    children?: Directory[];
    playlists: Playlist[];
  }

  export interface Link extends Base {
    source_directory_id?: number;
    source_playlist_id?: number;
    target_directory_id?: number;
    target_playlist_id?: number;
  }

  export interface Task {
    uid: string;
    name: string;
    status: string;
    next_run?: string;
    last_status?: string;
    last_run?: string;
    last_message?: string;
    last_error?: string;
    interval?: number;
    recurring: boolean;
  }

  export interface TaskHistory extends Base {
    name: string;
    result: string;
    run_at: string;
    message?: string;
    error?: string;
    duration: number;
  }

  export interface Track extends Base {
    spotify_id: string;
    name: string;
  }

  export interface TrackHistory extends Track {
    history_id: number;
    played_at: string;
    play_count?: number;
  }

  export interface TrackAdded extends Track {
    playlist: Playlist;
    created_at: string;
  }

  export interface TrackDeleted extends Track {
    playlist: Playlist;
    deleted_at: string;
  }

  export interface TrackDuplicate extends Track {
    amount: number;
  }

  export interface GeneratorWindow {
    start: string;
    end: string;
    min_plays: number;
    burst_interval_s: number;
  }

  export interface Generator extends Base {
    name: string;
    description?: string;
    params: {
      track_amount: number;
      excluded_playlist_ids?: number[];
      excluded_track_ids?: number[];
      preset: string;
      params_custom?: {};
      params_forgotten?: {};
      params_top?: {
        window: GeneratorWindow;
      };
      params_old_top?: {
        peak_window: GeneratorWindow;
        current_window: GeneratorWindow;
      };
    }
  }
}
