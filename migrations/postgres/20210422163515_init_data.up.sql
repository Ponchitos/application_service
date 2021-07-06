insert into public.installation_types (name) values ('AVAILABLE'), ('PREINSTALLED'), ('FORCE_INSTALLED'), ('REQUIRED_FOR_SETUP');
insert into public.location_types (name) values ('GOOGLE_PLAY'), ('SELF_HOSTED');

insert into public.installation_type_location_type_set (location_id, type_id)
    values (
               (select location_id from public.location_types where name = 'GOOGLE_PLAY'),
               (select type_id from public.installation_types where name = 'AVAILABLE')
           ),
           (
               (select location_id from public.location_types where name = 'GOOGLE_PLAY'),
               (select type_id from public.installation_types where name = 'PREINSTALLED')
           ),
           (
               (select location_id from public.location_types where name = 'GOOGLE_PLAY'),
               (select type_id from public.installation_types where name = 'FORCE_INSTALLED')
           ),
           (
               (select location_id from public.location_types where name = 'GOOGLE_PLAY'),
               (select type_id from public.installation_types where name = 'REQUIRED_FOR_SETUP')
           ),
           (
               (select location_id from public.location_types where name = 'SELF_HOSTED'),
               (select type_id from public.installation_types where name = 'PREINSTALLED')
           ),
           (
               (select location_id from public.location_types where name = 'SELF_HOSTED'),
               (select type_id from public.installation_types where name = 'FORCE_INSTALLED')
           ),
           (
               (select location_id from public.location_types where name = 'SELF_HOSTED'),
               (select type_id from public.installation_types where name = 'REQUIRED_FOR_SETUP')
           );