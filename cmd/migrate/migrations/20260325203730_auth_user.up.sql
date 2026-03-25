create schema if not exists auth;

create table if not exists auth.users (
  id uuid primary key,
  email text not null unique,
  password_hash text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists public.profile (
  user_id uuid primary key references auth.users(id) on delete cascade,
  firstname text,
  surname text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
