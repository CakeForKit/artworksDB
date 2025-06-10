CREATE TABLE IF NOT EXISTS artworks.migrations (
    version UInt32,
    name String,
    applied_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY version;


select * from artworks."Artworks"