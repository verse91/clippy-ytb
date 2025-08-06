-- Enable pgcrypto extension for gen_random_uuid()
create extension if not exists pgcrypto;

create table if not exists downloads (
  id uuid primary key default gen_random_uuid(),
  url text not null,
  status text,
  message text,
  created_at timestamptz default now()
);

create table if not exists time_range_downloads (
  id uuid primary key default gen_random_uuid(),
  url text not null,
  start_time integer not null,
  end_time integer not null,
  status text default 'completed',
  message text,
  output_file text,
  created_at timestamptz default now()
);

-- Create profiles table to store user credits and email from auth.users
create table profiles (
  id uuid primary key references auth.users(id) on delete cascade,
  email text unique,
  credits integer default 0 check (credits >= 0),
  -- additional optional fields can be added here as needed
  created_at timestamptz default now()
);

-- Enable Row Level Security
alter table profiles enable row level security;

-- Create policies
create policy "Users can view own profile" on profiles for select using (auth.uid() = id);
create policy "Users can update own profile" on profiles for update using (auth.uid() = id)
  with check (auth.uid() = id and credits = (select credits from profiles where id = auth.uid()));

-- Create function to handle new user signup, copying email from auth.users
create or replace function public.handle_new_user()
returns trigger as $$
begin
  insert into public.profiles (id, email, credits)
  values (new.id, new.email, 0);
  return new;
end;
$$ language plpgsql security definer set search_path = public, auth;

-- Create trigger for new user signup
create or replace trigger on_auth_user_created
  after insert on auth.users
  for each row execute procedure public.handle_new_user();
