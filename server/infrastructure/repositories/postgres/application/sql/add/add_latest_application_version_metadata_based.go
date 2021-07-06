package add

const LatestApplicationVersionMetadataBasedSQL = `
	insert into public.application_latest_versions
		(application_id, version_id, enterprise_id)
	values 
		($1, $2, $3);
`
