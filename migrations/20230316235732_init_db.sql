-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons (
    id uuid PRIMARY KEY NOT NULL,
    email text,
    telegram_id text
);

CREATE TYPE manga_source_type AS ENUM ('mangalib', 'mangadex');
CREATE TABLE manga_info (
    id uuid PRIMARY KEY NOT NULL,
    manga_name text NOT NULL,
    url text NOT NULL,
    source_type manga_source_type NOT NULL,
    last_updated_number int
);

CREATE TABLE subscribers (
    id int PRIMARY KEY NOT NULL,
    person_id uuid NOT NULL,
    manga_id uuid NOT NULL,
    CONSTRAINT person_id_constr FOREIGN KEY(person_id) REFERENCES persons(id),
    CONSTRAINT manga_id_constr FOREIGN KEY(manga_id) REFERENCES manga_info(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE persons;
DROP TABLE manga_info;
DROP TABLE subscribers;
DROP TYPE manga_source_type;
-- +goose StatementEnd
