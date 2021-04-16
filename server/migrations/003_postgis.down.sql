ALTER TABLE walk
    DROP COLUMN start_point,
    DROP COLUMN end_point;

DROP EXTENSION IF EXISTS cube;
DROP EXTENSION IF EXISTS earthdistance;
