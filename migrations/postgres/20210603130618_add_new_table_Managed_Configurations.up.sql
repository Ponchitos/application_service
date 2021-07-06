create table if not exists public.managed_configurations (
    managed_configuration_id    bigserial primary key,

    uuid                        uuid default uuid_generate_v4() not null,

    group_uuid                  uuid not null,

    configuration               bytea not null,

    application_version_id      integer not null,

    created                     timestamp without time zone default now(),
    modified                    timestamp without time zone,
    deleted                     timestamp without time zone,

    constraint managed_configurations__application_version_id foreign key (application_version_id) references public.application_versions (version_id) on delete cascade
);
create index if not exists idx_managed_configurations__group_uuid on public.managed_configurations (group_uuid);
create index if not exists idx_managed_configurations__uuid on public.managed_configurations (uuid);
create unique index if not exists uk_managed_configurations__application_version_id on public.managed_configurations (application_version_id);