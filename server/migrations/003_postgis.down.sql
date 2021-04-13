DROP EXTENSION postgis;

ALTER TABLE walk
    DROP COLUMN start_point,
    DROP COLUMN end_point;
