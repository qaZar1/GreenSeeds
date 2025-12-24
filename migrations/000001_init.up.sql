CREATE SCHEMA IF NOT EXISTS green_seeds;

CREATE TABLE IF NOT EXISTS green_seeds.bunkers (
    bunker INT,
    distance INT,
    PRIMARY KEY (bunker)
);

CREATE TABLE IF NOT EXISTS green_seeds.seeds (
    seed_ru VARCHAR(50),
    seed VARCHAR(50),
    min_density INT,
    max_density INT,
    tank_capacity INT,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (seed)
);

CREATE TABLE IF NOT EXISTS green_seeds.placement (
    bunker INT UNIQUE,
    seed VARCHAR(50),
    amount INT DEFAULT 0,
    FOREIGN KEY (bunker) REFERENCES green_seeds.bunkers(bunker),
    FOREIGN KEY (seed) REFERENCES green_seeds.seeds(seed)
);

CREATE TABLE IF NOT EXISTS green_seeds.receipts (
    receipt SERIAL,
    seed VARCHAR(50),
    gcode TEXT,
    updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (seed) REFERENCES green_seeds.seeds(seed),
    PRIMARY KEY (receipt)
);

CREATE TABLE IF NOT EXISTS green_seeds.assignments (
    id SERIAL,
    shift BIGINT,
    number INT,
    receipt BIGINT,
    amount INT,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (receipt) REFERENCES green_seeds.receipts(receipt),
    FOREIGN KEY (shift) REFERENCES green_seeds.shifts(shift),
    PRIMARY KEY (shift, number, receipt)
);

CREATE TABLE IF NOT EXISTS green_seeds.shifts (
    shift SERIAL,
    dt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(20),
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (username) REFERENCES green_seeds.users(username),
    PRIMARY KEY (shift)
);

CREATE TABLE IF NOT EXISTS green_seeds.users (
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name DVARCHAR(50) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS green_seeds.reports (
    id SERIAL,
    shift BIGINT,
    number INT,
    receipt BIGINT,
    turn INT,
    dt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN,
    error VARCHAR(255),
    solution VARCHAR(50),
    mark VARCHAR(50),
    FOREIGN KEY (shift, number, receipt)
        REFERENCES green_seeds.assignments (shift, number, receipt),
    PRIMARY KEY (shift, number, receipt, turn)
);

CREATE TABLE IF NOT EXISTS green_seeds.logs (
    id SERIAL,
    dt TIMESTAMP WITH TIME ZONE,
    lvl VARCHAR(10),
    request_id VARCHAR(255),
    msg TEXT,
    caller VARCHAR(512),
    username VARCHAR(20),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS green_seeds.device_settings (
    key VARCHAR(30),
    value TEXT
)




------------- FOR SQLITE -------------
CREATE TABLE IF NOT EXISTS calibration (
    session_id TEXT PRIMARY KEY,
    first_photo_path TEXT,
    second_photo_path TEXT,
    dx REAL,
    dy REAL,
    cir REAL,
    d_per_step REAL,
    created_at TEXT DEFAULT (datetime('now'))
);


