package delete

const LatestVersionSQL = `
	update public.application_latest_versions
		set deleted = now()
	where application_id = (select application_id from public.applications where uuid = $1 and enterprise_id = $2 and deleted is null) and
		  deleted is null;
`
