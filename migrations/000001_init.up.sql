CREATE SCHEMA IF NOT EXISTS green_seeds;

CREATE TABLE IF NOT EXISTS green_seeds.bunkers (
    bunker INT,
    distance INT,
    PRIMARY KEY (bunker)
);

CREATE TABLE IF NOT EXISTS green_seeds.seeds (
    seed VARCHAR(50),
    min_density INT,
    max_density INT,
    tank_capacity INT,
    latency INT,
    PRIMARY KEY (seed)
);

CREATE TABLE IF NOT EXISTS green_seeds.placement (
    bunker INT,
    seed VARCHAR(50),
    FOREIGN KEY (bunker) REFERENCES green_seeds.bunkers(bunker),
    FOREIGN KEY (seed) REFERENCES green_seeds.seeds(seed)
);

CREATE TABLE IF NOT EXISTS green_seeds.receipts (
    receipt SERIAL,
    seed VARCHAR(50),
    gcode TEXT,
    updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    FOREIGN KEY (seed) REFERENCES green_seeds.seeds(seed),
    PRIMARY KEY (receipt)
);

CREATE TABLE IF NOT EXISTS green_seeds.assignments (
    shift BIGINT,
    number INT,
    receipt BIGINT,
    amount INT,
    FOREIGN KEY (receipt) REFERENCES green_seeds.receipts(receipt),
    FOREIGN KEY (shift) REFERENCES green_seeds.shifts(shift),
    PRIMARY KEY (shift, number, receipt)
);

CREATE TABLE IF NOT EXISTS green_seeds.shifts (
    shift SERIAL,
    dt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(20),
    FOREIGN KEY (username) REFERENCES green_seeds.users(username),
    PRIMARY KEY (shift)
);

CREATE TABLE IF NOT EXISTS green_seeds.users (
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name DVARCHAR(50) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS green_seeds.reports (
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