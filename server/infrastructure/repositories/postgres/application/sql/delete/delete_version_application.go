package delete

const VersionApplicationSQL = `
	update public.application_versions
		set deleted = now()
	where uuid = $1 and
		  application_id = (select application_id from public.applications where uuid = $2 and enterprise_id = $3 and deleted is null)
	returning application_metadata_id;
`
