-- https://hashrocket.com/blog/posts/juxtaposing-earthdistance-and-postgis
-- -> no postgis is simpler and we don't need to be extra accurate
-- -> "cube" distance is best:
--    SELECT *, earth_distance(ll_to_earth(point[0], point[1]),  ll_to_earth(32, 2)) as d from table;
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

ALTER TABLE walk
    ADD COLUMN start_point POINT NOT NULL,
    ADD COLUMN end_point   POINT NOT NULL;
