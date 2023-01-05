-- explicit type casting is required because of how statements are prepared on the pgserver

--name: get-search-bounds
-- Computes optimum bounds for series generation.
-- $1: search_batch_size
-- $2: head_cursor
-- $3: head_block_lag
WITH low AS (
	SELECT lower_bound FROM syncer_meta
)
SELECT low.lower_bound, LEAST(low.lower_bound::int+$1::int, $2::int -$3::int) AS upper_bound FROM low

--name: get-missing-blocks
-- Generates a bounded series and searches for missing blocks.
-- $1: series_lower_bound
-- $2: series_upper_bound
SELECT s.i AS missing_blocks
FROM generate_series($1::int, $2::int) s(i)
WHERE NOT EXISTS (SELECT 1 FROM blocks WHERE n = s.i)

--name: set-search-lower-bound
-- Updates the search lower bound to keep the query efficient during janitor sweeps.
-- $1: new_lower_bound
UPDATE syncer_meta SET lower_bound = $1::int

--name: init-syncer-meta
-- If first run, populate syncer_meta table. Safe to run multiple times.
-- $1: init_lower_bound
INSERT INTO syncer_meta(lower_bound) SELECT $1::int WHERE NOT EXISTS (
	SELECT lower_bound FROM syncer_meta
);

--name: commit-block
-- All transaction in this block have been processed through the pipleine without error.
-- $1: block_number
INSERT INTO blocks(n) VALUES($1::int)
