package edit

const RollbackUninstallApplicationFromGroupSQL = `
	update public.application_group_set
		set status = previous_status
	where application_version_id = (select version_id from public.application_versions where uuid = $1 and deleted is null) and
		  application_id = (select application_id from public.applications where uuid = $2 and deleted is null) and
		  group_uuid = $3 and
		  enterprise_id = $4 and
		  deleted is null;
`
