package get

const ApplicationToGroupSetsWithDeletedRowsSQL = `
	select set_id,
		   application_id,
		   application_version_id,
		   group_uuid,
		   group_name,
		   enterprise_id,
		   status,
		   previous_status,
		   created,
		   modified,
		   deleted	
	from public.application_group_set
	where application_id = (select application_id from public.applications where uuid = $1 and deleted is null) and
		  application_version_id = (select version_id from public.application_versions where uuid = $1 and deleted is null) and
		  group_uuid = $3 and
		  enterprise_id = $4;
`
