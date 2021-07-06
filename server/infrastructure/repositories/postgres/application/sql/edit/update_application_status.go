package edit

const ApplicationStatusByVersionUUIDSQL = `
	update public.applications
		set status = $1
	where application_id = (select application_id from public.application_versions where uuid = $2 and deleted is null) and
		  enterprise_id = $3 and
		  deleted is null;
`
