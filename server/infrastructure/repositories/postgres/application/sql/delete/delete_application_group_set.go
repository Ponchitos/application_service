package delete

const ApplicationGroupSetSQL = `
	update public.application_group_set
		set deleted = now(),
			status = $1
	where application_version_id = (select version_id from public.application_versions where uuid = $2 and deleted is null) and
		  application_id = (select application_id from public.applications where uuid = $3 and deleted is null) and	
		  enterprise_id = $4 and
		  group_uuid = $5;
`
