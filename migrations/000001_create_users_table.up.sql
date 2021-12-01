create table if not exists users(
id bigserial PRIMARY KEY,
email citext UNIQUE NOT NULL,
apiKey text not null,
apiCount integer not null default 10,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
version integer NOT NULL DEFAULT 1
);	