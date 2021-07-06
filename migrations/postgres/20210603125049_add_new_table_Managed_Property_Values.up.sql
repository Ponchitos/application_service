create table if not exists public.managed_property_values (
  managed_property_value_id bigserial primary key,

  group_uuid                uuid not null,

  managed_property_id       integer not null,

  value                     bytea,

  created                   timestamp without time zone default now(),
  modified                  timestamp without time zone,
  deleted                   timestamp without time zone,

  constraint  managed_property_values__managed_property_id foreign key (managed_property_id) references public.managed_properties (property_id) on delete cascade
);
create index if not exists idx_managed_property_values__group_uuid on public.managed_property_values (group_uuid);
create unique index if not exists uk_managed_property_values__managed_property_id__group_uuid on public.managed_property_values (managed_property_id, group_uuid);