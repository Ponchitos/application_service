package edit

const ApplicationToGroupStatusSQL = `
	update public.application_group_set
		set status = $1,
			modified = now()
	where application_version_id = (select version_id from public.application_versions where uuid = $2 and deleted is null) and
		  application_id = (select application_id from public.applications where uuid = $3 and deleted is null) and
		  group_uuid = $4 and
		  enterprise_id	= $5 and
		  deleted is null;	
`
