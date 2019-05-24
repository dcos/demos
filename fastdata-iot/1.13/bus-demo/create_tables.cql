CREATE KEYSPACE IF NOT EXISTS streaming WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

CREATE TABLE IF NOT EXISTS streaming.routeInfos (
    id text,
    display_name text,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS streaming.vehicles (
    id text,
    time timestamp,
    longitude double,
    latitude double,
    heading int,
    route_id text,
    run_id text,
    seconds_since_report int,
    PRIMARY KEY (id, time)
);

CREATE TABLE IF NOT EXISTS streaming.vehicles_flink (
    id text,
    time timestamp,
    longitude double,
    latitude double,
    heading int,
    route_id text,
    run_id text,
    seconds_since_report int,
    PRIMARY KEY (id, time)
);

CREATE TABLE IF NOT EXISTS streaming.routes (
    id text,
    order_id int,
    route_id text,
    longitude double,
    latitude double,
    display_name text,
    PRIMARY KEY (route_id, order_id, id)
) WITH CLUSTERING ORDER BY (order_id ASC);

CREATE TABLE IF NOT EXISTS streaming.vehicles_by_tileid (
    tile_id text,
    time_id bigint,
    id text,
    time timestamp,
    longitude double,
    latitude double,
    heading int,
    route_id text,
    run_id text,
    seconds_since_report int,
    PRIMARY KEY ((tile_id, time_id), time, id)
) WITH CLUSTERING ORDER BY (time DESC);

CREATE TABLE IF NOT EXISTS streaming.vehicles_by_tileid_flink (
  tile_id text,
  time_id bigint,
  id text,
  time timestamp,
  longitude double,
  latitude double,
  heading int,
  route_id text,
  run_id text,
  seconds_since_report int,
  PRIMARY KEY ((tile_id, time_id), time, id)
) WITH CLUSTERING ORDER BY (time DESC);

CREATE TABLE IF NOT EXISTS streaming.vehiclecluster (
    id int,
    time_stamp bigint,
    longitude double,
    latitude double,
    amount int,
    PRIMARY KEY (id, time_stamp)
);

CREATE TABLE IF NOT EXISTS streaming.vehiclecluster_flink (
  id int,
  time_stamp bigint,
  longitude double,
  latitude double,
  amount int,
  PRIMARY KEY (id, time_stamp)
);

CREATE TABLE IF NOT EXISTS streaming.vehiclecluster_by_tileid (
    tile_id text,
    id int,
    time_stamp bigint,
    longitude double,
    latitude double,
    amount int,
    PRIMARY KEY (tile_id, id, time_stamp)
);

CREATE TABLE IF NOT EXISTS streaming.vehiclecluster_by_tileid_flink (
  tile_id text,
  id int,
  time_stamp bigint,
  longitude double,
  latitude double,
  amount int,
  PRIMARY KEY (tile_id, id, time_stamp)
);

CREATE TABLE IF NOT EXISTS streaming.vehicleclusterdetails (
    id bigint,
    pos_id int,
    time_stamp bigint,
    longitude double,
    latitude double,
    PRIMARY KEY (id, pos_id, time_stamp)
);

CREATE TABLE IF NOT EXISTS streaming.vehicleclusterdetails_flink (
  id bigint,
  pos_id int,
  time_stamp bigint,
  longitude double,
  latitude double,
  PRIMARY KEY (id, pos_id, time_stamp)
);