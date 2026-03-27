alter table auth.session
add column if not exists csrf_token text not null default '';
