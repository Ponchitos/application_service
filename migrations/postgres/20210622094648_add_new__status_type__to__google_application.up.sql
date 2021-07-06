create type public.google_application_status as enum ('APPROVED', 'UNAPPROVED');

alter table if exists public.google_application_info add column status google_application_status default 'APPROVED';