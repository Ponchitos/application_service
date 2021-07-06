create extension "uuid-ossp";

create type public.location_type as enum ('SELF_HOSTED', 'GOOGLE_PLAY');
create type public.available_type as enum ('PRIVATE', 'PUBLIC');
create type public.application_status as enum ('APPROVED', 'INSTALLED', 'UPDATE_AVAILABLE', 'WAITING_INSTALL', 'WAITING_UPDATE', 'WAITING_UNINSTALL', 'UNINSTALLED');

create table if not exists public.application_metadata (
    metadata_id             bigserial primary key,

    uuid                    uuid default uuid_generate_v4() not null,

    link                    varchar not null,

    package_name            varchar not null,
    application_label       varchar not null,
    version_name            varchar not null,
    file_size               varchar not null,
    file_sha1_base64        varchar not null,
    file_sha256_base64      varchar not null,
    icon_base64             varchar not null,
    externally_hosted_url   varchar not null,

    native_codes            varchar[],
    certificate_base64s     varchar[],
    uses_features           varchar[],

    uses_permissions        bytea,

    version_code            integer not null,
    maximum_sdk             integer not null,

    created                 timestamp without time zone default now(),
    deleted                 timestamp without time zone
);
create index if not exists idx_application_metadata__uuid on public.application_metadata (uuid);


create table if not exists public.google_application_info (
    id          bigserial primary key,

    uuid        uuid default uuid_generate_v4() not null,

    name        varchar not null,
    title       varchar not null,

    created     timestamp without time zone default now(),
    deleted     timestamp without time zone
);
create index if not exists idx_google_application_info__uuid on public.google_application_info (uuid);
create index if not exists idx_google_application_info__name on public.google_application_info (name);


create table if not exists public.application_permissions (
    id                          bigserial primary key,

    google_application_info_id  integer not null,

    permission_id               varchar not null,
    name                        varchar not null,
    description                 varchar not null,

    created                     timestamp without time zone default now(),
    deleted                     timestamp without time zone,

    constraint application_permissions__google_application_info_id foreign key (google_application_info_id) references public.google_application_info (id) on delete cascade
);
create index if not exists idx_application_permissions__google_application_info_id on public.application_permissions (google_application_info_id);


create table if not exists public.application_tracks (
    id                            bigserial primary key,

    google_application_info_id    integer not null,

    track_id                      varchar not null,
    track_alias                   varchar not null,

    created                       timestamp without time zone default now(),
    deleted                       timestamp without time zone,

    constraint application_tracks__google_application_info_id foreign key (google_application_info_id) references public.google_application_info (id) on delete cascade
);
create index if not exists idx_application_tracks__google_application_info_id on public.application_tracks (google_application_info_id);


create table if not exists public.managed_properties (
    property_id                 bigserial primary key,
    google_application_info_id  integer,

    key                         varchar not null,
    type                        varchar not null,
    title                       varchar not null,
    description                 varchar not null,

    default_value               jsonb,

    entries                     bytea,

    created                     timestamp without time zone default now(),
    deleted                     timestamp without time zone,

    constraint managed_properties__google_application_info_id foreign key (google_application_info_id) references public.google_application_info (id) on delete cascade
);


create table if not exists public.managed_properties_set (
    set_id                      bigserial primary key,

    parent_managed_property_id  integer not null,
    child_managed_property_id   integer not null,

    created                     timestamp without time zone default now(),
    deleted                     timestamp without time zone,

    constraint managed_properties_set__parent_managed_property_id foreign key (parent_managed_property_id) references public.managed_properties (property_id) on delete cascade,
    constraint managed_properties_set__child_managed_property_id foreign key (child_managed_property_id) references public.managed_properties (property_id) on delete cascade
);
create index if not exists idx_managed_properties_set__parent_managed_property_id on public.managed_properties_set (parent_managed_property_id);
create index if not exists idx_managed_properties_set__child_managed_property_id on public.managed_properties_set (child_managed_property_id);

create table if not exists public.applications (
    application_id   bigserial primary key,

    uuid             uuid default uuid_generate_v4() not null,

    package_name     varchar not null,

    name             varchar not null,

    enterprise_id    varchar not null,

    available        available_type not null,
    location         location_type not null,

    status           application_status not null,

    created          timestamp without time zone default now(),
    deleted          timestamp without time zone
);
create index if not exists idx_applications__uuid on public.applications (uuid);
create index if not exists idx_applications__enterprise_id on public.applications (enterprise_id);
create index if not exists idx_applications__package_name on public.applications (package_name);


create table if not exists public.application_versions (
    version_id                  bigserial primary key,

    uuid                        uuid default uuid_generate_v4() not null,

    application_id              integer not null,

    application_metadata_id     integer,

    google_application_info_id  integer,

    version_code                integer,
    min_sdk                     integer,

    version_name                varchar,
    icon                        varchar,

    created                     timestamp without time zone default now(),
    deleted                     timestamp without time zone,

    constraint application_versions__application_id foreign key (application_id) references public.applications (application_id) on delete cascade,
    constraint application_versions__application_metadata_id foreign key (application_metadata_id) references public.application_metadata (metadata_id) on delete cascade,
    constraint application_versions__google_application_info_id foreign key (google_application_info_id) references public.google_application_info (id) on delete cascade
);
create index if not exists idx_application_versions__uuid on public.application_versions (uuid);
create index if not exists idx_application_versions__application_id on public.application_versions (application_id);
create index if not exists idx_application_versions__application_metadata_id on public.application_versions (application_metadata_id);
create index if not exists idx_application_versions__google_application_info_id on public.application_versions (google_application_info_id);


create table if not exists public.application_latest_versions (
    latest_version_id   bigserial primary key,

    application_id      integer not null,

    version_id          integer not null,

    enterprise_id       varchar not null,

    created             timestamp without time zone default now(),
    deleted             timestamp without time zone,
    modified            timestamp without time zone,

    constraint application_latest_versions__application_id foreign key (application_id) references public.applications (application_id) on delete cascade,
    constraint application_latest_versions__version_id foreign key (version_id) references public.application_versions (version_id) on delete cascade
);
create index if not exists idx_application_latest_versions__application_id on public.application_latest_versions (application_id);
create index if not exists idx_application_latest_versions__version_id on public.application_latest_versions (version_id);
create index if not exists idx_application_latest_versions__enterprise_id on public.application_latest_versions (enterprise_id);


create table if not exists public.installation_types (
    type_id     bigserial primary key,

    name        varchar not null,

    is_active   boolean default true,

    created     timestamp without time zone default now(),
    deleted     timestamp without time zone
);
create index if not exists idx_installation_types__name on public.installation_types (name);


create table if not exists public.location_types (
    location_id     bigserial primary key,

    name            varchar not null,

    is_active       boolean default true,

    created         timestamp without time zone default now(),
    deleted         timestamp without time zone
);
create index if not exists idx_location_types__name on public.location_types (name);


create table if not exists public.installation_type_location_type_set (
    set_id      bigserial primary key,

    location_id integer not null,
    type_id     integer not null,

    is_active   boolean default true,

    created     timestamp without time zone default now(),
    deleted     timestamp without time zone,

    constraint installation_type_location_type_set__location_id foreign key (location_id) references public.location_types (location_id) on delete cascade,
    constraint installation_type_location_type_set__type_id foreign key (type_id) references public.installation_types (type_id) on delete cascade
);

create index if not exists idx_installation_type_location_type_set__location_id on public.installation_type_location_type_set (location_id);
create index if not exists idx_installation_type_location_type_set__type_id on public.installation_type_location_type_set (type_id);