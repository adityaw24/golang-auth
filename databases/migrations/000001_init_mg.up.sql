begin;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

Create table if not exists public.users (
  user_id uuid primary key default uuid_generate_v4(),
  username varchar(200) not null unique,
  name TEXT not null,
  password varchar(200) not null,
  email varchar(200) not null unique,
  phone varchar(200) not null unique,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp,
  is_delete boolean default false,

  constraint username_unique unique (username),
  constraint email_unique unique (email),
  constraint phone_unique unique (phone)
);

commit;