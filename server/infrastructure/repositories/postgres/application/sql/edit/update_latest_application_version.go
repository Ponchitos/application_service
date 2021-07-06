package edit

const LatestApplicationSQL = `
	update public.application_latest_versions
		set version_id = $1,
			modified = now()
	where application_id = (select application_id from public.applications where uuid = $2 and deleted is null) and
		  enterprise_id = $3 and
		  deleted is null;
`
