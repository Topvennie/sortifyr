-- name: GeneratorTrackCreateBatch :exec
INSERT INTO generator_tracks (generator_id, track_id)
VALUES (
  UNNEST($1::int[]),
  UNNEST($2::int[])
);

-- name: GeneratorTrackDeleteByGenerator :exec
DELETE FROM generator_tracks
WHERE generator_id = $1;
