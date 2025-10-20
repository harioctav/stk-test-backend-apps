-- Add UUID column (nullable first for existing records)
ALTER TABLE menus ADD COLUMN uuid VARCHAR(36) UNIQUE AFTER id;

-- Generate UUIDs for existing records
UPDATE menus SET uuid = UUID() WHERE uuid IS NULL;

-- Make it NOT NULL after population
ALTER TABLE menus MODIFY COLUMN uuid VARCHAR(36) UNIQUE NOT NULL;

-- Add index for faster lookup
CREATE INDEX idx_uuid ON menus(uuid);

