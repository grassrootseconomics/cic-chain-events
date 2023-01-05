CREATE TABLE IF NOT EXISTS blocks (
	n INT UNIQUE PRIMARY KEY
);

CREATE TABLE syncer_meta (
  lower_bound INT
);

---- create above / drop below ----

DROP TABLE syncer_meta;
DROP TABLE blocks;
