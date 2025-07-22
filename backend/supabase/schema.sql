create table if not exists downloads (
  id uuid primary key default gen_random_uuid(),
  url text not null,
  status text,
  message text,
  created_at timestamp default now()
);
