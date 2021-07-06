package delete

const ApplicationMetadataByUUIDSQL = `
	update public.application_metadata
		set deleted = now()
	where uuid = $1;
`

const ApplicationMetadataByIDSQL = `
	update public.application_metadata
		set deleted = now()
	where metadata_id = $1
	returning uuid;
`
