CREATE EXTENSION postgis;

ALTER TABLE walk
    ADD COLUMN start_point geography(point),
    ADD COLUMN end_point   geography(point);
