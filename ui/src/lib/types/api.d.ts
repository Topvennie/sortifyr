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
    public: boolean;
    track_amount: number;
    collaborative: boolean;
    has_cover: boolean;
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
    next_run: string;
    last_status?: string;
    last_run?: string;
    last_message?: string;
    last_error?: string;
    interval?: number;
  }

  export interface TaskHistory extends Base {
    name: string;
    result: string;
    run_at: string;
    error?: string;
    duration: number;
  }
}
