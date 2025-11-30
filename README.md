# Spotify organizer

A web application to help organize and automate your Spotify playlists.

While this tool is fully functional for general use, it is designed around __my own playlist organization system__.
So a few features may seem odd at first glance.

Nevertheless, many of the concepts are broadly useful for anyone with a large / structured Spotify library.

To give an idea for whom this might by useful, this is my setup:

<details>
  <summary>Personal setup</summary>

  I organize playlists into directory trees.
  My two main root directories are:

  __All__: A complete archive of everything \
  __Good__: A curated version I listen to daily

  These directories mirror each other for the most part in structure.
  
  When I add new music, I add it to a playlist under __Good__.
  This tool automatically copies it to the correct playlist(s) in the __All__ directory.
  Later, when I grow tired of the track, I remove it from __Good__, but it remains archived under __All__ so I never lose it.

  A simplified schematic:

```
Root/
  All/
    Genres/
      Pop
      Rock
    Instrument/
      Piano
      Guitar
  Good/
    Genres/
      Pop
      Rock
    Instrument/
      Piano
      Guitar
```

</details>

## Features

__Playlists__

- Fetches all your added playlists and their tracks (also playlists you saved from someone else).
- Spotify does __not__ expose Spotify-created playlists that you saved, so those cannot be fetched.

__Directories__

Spotify does not expose playlist folders through the API, that's why this project includes its own internal directory system.

It mirrors the Spotify system but you have to manually keep it in sync.

While it is not necessary for other features that you use the directory system, it is recommended.

__Links__

A link defines a one-way relationship between a _source_ and a _target_:

- The _source_ can be a playlist or a directory
- The _target_ can be a playlist or a directory
- Tracks from the source are automatically copied into the target

Examples:

- Source: a directory \
  -> Target: a playlist \
  -> Result: The playlist will contain all tracks from all playlists under the directory (including nested subdirectories).
- Source: a playlist \
  -> Target: a directory \
  -> Result: all playlists inside that directory (or nested subdirectories) receive the new tracks.

This system enables setups like:

- Always have a _master playlist_ containing everything I saved.
- Mirror playlists into archive playlists.
- Build automatic multi genre aggregators.

### Future features

Have a look at the issues to get a sneak peak of upcoming features!

## Spotify App setup

To use this tool you must create a Spotify Developer App, which provides the credentials needed for authentication.

1. Go to the [Spotify developers website](https://developer.spotify.com/).
2. Open your profile menu -> __Dashboard__.
3. Click __Create App__.
4. Fill out the form. Pay attention to:
    - __Redirect URIs__:
        - Reverse Proxy: `https://<domain>/api/auth/callback/spotify`
        - Local: `http://127.0.0.1:<port>/api/auth/callback/spotify` (default production port: __8000__, development: __3001__)
    - __Which API are you planning to use?__:
        - Choose Web API
5. After creation you can access your __client ID__ and __client secret__.

## Production Deployment

__Important Notes__

- The webapp uses Spotify authentication, meaning _anyone with a Spotify account can log into your instance_.
  The application supports multiple users, but all actions will occur through __your__ Spotify Developer App credentials.
  Worst-case scenario: someone could intentionally overload your app with API calls, causing temporary rate limiting.
- Spotify has clear [terms and conditions](https://developer.spotify.com/terms) for the usage of their web API.
  If you believe this project violates any policies, please open an issue so it can be corrected.

### Recommended Deployment (Docker)

1. Copy `docker-compose.prod.yml` -> `docker-compose.yml`.
2. Copy `.env.prod.example` -> `.env`.
3. Create your Spotify app and fill in the environment variables.
4. Run `docker compose up`.
5. The server is reachable on port __8000__.

To update:

```bash
docker compose pull
docker compose down
docker compose up
```

### Manual Setup (Advanced)

The container image is published as: `ghcr.io/topvennie/spotify_organizer`.

Required additional services:

- __Postgres__
- __Redis__
- __S3 compatible storage__

Configuration variables can be overridden via environment variables.
Keys use uppercase and replace `.` with `_` (e.g. `redis.url` -> `REDIS_URL`)

Search the codebase for `config.Get` to view available configuration settings.
You will probably need to change a couple to use your own services.

## Development

> [!IMPORTANT]
> Due to Spotify redirect restrictions, authentication only works when visiting the website through [127.0.0.1:3000](http://127.0.0.1:3000).
> Do not use localhost as it will fail to authenticate.

### Quick Start

1. Install the tools listed in the [asdf file](./.tool-versions) (if you have _asdf_ run `asdf install`).
2. Install _make_.
3. Run `make setup` to install:
    - Backend tools: _Air_, _Goose_, _Sqlc_, _Deadcode_
    - Frontend depedencies
4. Install the git hooks for code quality: `git config --local core.hooksPath .githooks/`.
5. Copy `.env.example` -> `.env`.
6. Create a Spotify app (see above).
7. Fill in all `.env` values.
8. Run database migrations: `make migrate`.
9. Start the project: `make watch`.

Endpoints:

- __Backend__: <http://127.0.0.1:3001>
- __Frontend__: <http://127.0.0.1:3000>

### Makefile Commands

A makefile is used to simplify some tasks.
For an overview of all commands see the makefile.

A few common commands:

__Start the full stack__

```bash
make watch
```

Starts backend + frontend with hot module reloading.
(Requires restart after adding or removing dependencies).

__Create a new migration__

```bash
make create-migration
```

Prompts fro a name and then creates a new migration under `db/migrations`.
Edit the SQL, optionally add new queries under `db/queries`, then run:

```bash
make query
```

__Update SQLC queries__

```bash
make query
```

Parses migrations and queries to generated types SQL query code.
The result can be found in `pkg/sqlc`.
