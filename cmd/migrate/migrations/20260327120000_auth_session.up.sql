create table if not exists auth.session (
  session_token text primary key,
  user_id bigint not null references auth.users(id) on delete cascade,
  created_at timestamptz not null default now()
);
