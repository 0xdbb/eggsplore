--
-- Name: concessions; Type: TABLE; Schema:  Owner: -
--

CREATE TABLE concessions (
    name text,
    owner text,
    type text,
    status text,
    start_date timestamp without time zone,
    expiry_dat timestamp without time zone,
    assets text,
    pid integer,
    geom geometry(Polygon,4326)
);


--
-- Name: districts; Type: TABLE; Schema:  Owner: -
--

CREATE TABLE districts (
    region text,
    district text,
    geom geometry(Polygon,4326)
);


--
-- Name: forest_reserves; Type: TABLE; Schema:  Owner: -
--

CREATE TABLE forest_reserves (
    name_1 text,
    category text,
    name text,
    geom geometry(Polygon,4326)
);


--
-- Name: mining_static; Type: TABLE; Schema:  Owner: -
--

CREATE TABLE mining_static (
    id                  text PRIMARY KEY,
    geom                geometry(Polygon, 4326),
    area                double precision,
    status              varchar(255) DEFAULT 'OPEN',
    severity            text,
    severity_type       text,
    severity_score      bigint,
    proximity_to_water  boolean,
    inside_forest_reserve boolean,
    detection_date      timestamp without time zone,
    image_url           text,
    all_violation_types text,
    district            text,
    task_id             uuid NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    distance_to_water_m double precision,
    distance_to_forest_m double precision
);


--
-- Name: rivers; Type: TABLE; Schema:  Owner: -
--

CREATE TABLE rivers (
    id bigint,
    orig_fid bigint,
    geom geometry(Polygon,4326)
);


--
-- Name: idx_concessions_geom; Type: INDEX; Schema:  Owner: -
--

CREATE INDEX idx_concessions_geom ON concessions USING gist (geom);


--
-- Name: idx_districts_geom; Type: INDEX; Schema:  Owner: -
--

CREATE INDEX idx_districts_geom ON districts USING gist (geom);


--
-- Name: idx_forest_reserves_geom; Type: INDEX; Schema:  Owner: -
--

CREATE INDEX idx_forest_reserves_geom ON forest_reserves USING gist (geom);


--
-- Name: idx_mining_static_geom; Type: INDEX; Schema:  Owner: -
--

CREATE INDEX idx_mining_static_geom ON mining_static USING gist (geom);


--
-- Name: idx_rivers_geom; Type: INDEX; Schema:  Owner: -
--

CREATE INDEX idx_rivers_geom ON rivers USING gist (geom);


--
-- PostgreSQL database dump complete
--



