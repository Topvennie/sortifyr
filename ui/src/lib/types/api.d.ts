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

  export interface Setting extends Base {
    last_update?: string;
  }

  export interface Playlist extends Base {
    spotify_id: string;
    owner?: User;
    name: string;
    description?: string;
    public: boolean;
    tracks: number;
    collaborative: boolean;
  }

  export interface Directory extends Base {
    name: string;
    children?: Directory[];
    playlists: Playlist[];
  }
}
