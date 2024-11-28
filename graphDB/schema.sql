--Paths
CREATE TABLE graph_paths (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state INTEGER NOT NULL,
    size INTEGER NOT NULL,
    cars INTEGER NOT NULL
);

-- Path ends
CREATE TABLE path_ends (
    parent_id INTEGER PRIMARY KEY,
    longitude REAL NOT NULL,
    latitude REAL NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES graph_paths
)

--Nodes
CREATE TABLE graph_nodes (
    id INTEGER PRIMARY KEY UNIQUE,
    longitude REAL NOT NULL,
    latitude REAL NOT NULL
)