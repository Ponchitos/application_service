package add

const ApplicationGroupSetSQL = `
	insert into public.application_group_set
		(application_version_id, application_id, enterprise_id, group_name, group_uuid, status)
	values
		(
			(select version_id from public.application_versions where uuid = $1 and deleted is null),
			(select application_id from public.applications where uuid = $2 and deleted is null),
			$3,
			$4,
			$5,
			$6
		);
`
