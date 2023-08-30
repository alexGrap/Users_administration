CREATE TABLE IF NOT EXISTS segments
(
    id   SERIAL PRIMARY KEY ,
    segmentName TEXT NOT NULL,
    percents INT NOT NULL
);
CREATE TABLE IF NOT EXISTS subscription
(
    userId    BIGSERIAL,
    segment   SERIAL REFERENCES segments(id),
    PRIMARY KEY (userId, segment),
    timeToDie date NOT NULL
);