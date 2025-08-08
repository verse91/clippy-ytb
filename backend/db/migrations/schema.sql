-- Paste this into your Supabase sql editor

-- Enable pgcrypto extension for gen_random_uuid()
create extension if not exists pgcrypto;

-- Create schema_version table to track database schema versions
create table if not exists schema_version (
  id serial primary key,
  version integer not null,
  applied_at timestamptz default now()
);

-- Insert initial schema version
insert into schema_version (version) values (1) on conflict do nothing;

-- Create or replace the set_updated_at() trigger function
create or replace function public.set_updated_at()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
end;
$$ language plpgsql;

-- Drop any existing trigger on downloads table
drop trigger if exists set_updated_at_downloads on public.downloads;

create table if not exists public.downloads (
  id uuid primary key default gen_random_uuid(),
  url text not null,
  status text not null default 'pending' check (status in ('pending', 'processing', 'completed', 'failed')),
  message text,
  created_at timestamptz not null default now(),
  updated_at timestamptz default now()
);

-- Create trigger to automatically update updated_at on row updates
create trigger set_updated_at_downloads
  before update on public.downloads
  for each row execute procedure public.set_updated_at();

-- Drop any existing trigger on time_range_downloads table
drop trigger if exists set_updated_at_time_range_downloads on public.time_range_downloads;

create table if not exists public.time_range_downloads (
  id uuid primary key default gen_random_uuid(),
  url text not null,
  start_time integer not null,
  end_time integer not null check (end_time >= start_time),
  status text default 'pending',
  message text,
  output_file text,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Create trigger to automatically update updated_at on row updates
create trigger set_updated_at_time_range_downloads
  before update on public.time_range_downloads
  for each row execute procedure public.set_updated_at();

-- Create profiles table to store user credits and email from auth.users
create table if not exists profiles (
  id uuid primary key references auth.users(id) on delete cascade,
  email text unique,
  credits integer not null default 0 check (credits >= 0),
  -- additional optional fields can be added here as needed
  created_at timestamptz not null default now()
);

-- Enable Row Level Security
alter table profiles enable row level security;

-- Create policies with idempotent checks and proper role scoping
DO $$
BEGIN
  -- Drop existing policies if they exist
  IF EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'profiles' AND policyname = 'Users can view own profile') THEN
    DROP POLICY "Users can view own profile" ON profiles;
  END IF;

  IF EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'profiles' AND policyname = 'Users can update own profile') THEN
    DROP POLICY "Users can update own profile" ON profiles;
  END IF;

  -- Create policies with proper role scoping and NULL-safe comparison
  CREATE POLICY "Users can view own profile" ON profiles
    FOR SELECT TO authenticated
    USING (auth.uid() = id);

  CREATE POLICY "Users can update own profile" ON profiles
    FOR UPDATE TO authenticated
    USING (auth.uid() = id)
    WITH CHECK (auth.uid() = id AND credits IS NOT DISTINCT FROM (SELECT credits FROM profiles WHERE id = auth.uid()));
END $$;

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
