alter table auth.session
drop column if exists csrf_token;
