package add

const ApplicationVersionMetadataBasedSQL = `
	insert into public.application_versions
		(application_metadata_id, application_id, version_code, version_name, min_sdk, icon)
	values
		(
			$1,
			(select application_id from public.applications where uuid = $2 and deleted is null),
			$3,
			$4,
			$5,
			$6
		)
	returning version_id, uuid;
`
