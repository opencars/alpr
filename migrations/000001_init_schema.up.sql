CREATE EXTENSION "uuid-ossp";

CREATE TABLE "recognitions" (
	"id"         UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
	"image_key"  TEXT      NOT NULL,
	"plate"      TEXT      NOT NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX recognitions_plate_idx ON recognitions("plate");