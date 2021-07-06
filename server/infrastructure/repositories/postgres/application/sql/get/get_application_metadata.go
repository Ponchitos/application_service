package get

const ApplicationMetadataSQL = `
	select uuid,
		   metadata_id,
		   link,
		   package_name,
		   application_label,
		   version_name,
		   file_size,
		   file_sha1_base64,
		   file_sha256_base64,
		   icon_base64,
		   externally_hosted_url,
		   native_codes,
		   certificate_base64s,
		   uses_features,
		   version_code,
		   maximum_sdk,
		   created,
		   uses_permissions	
	from public.application_metadata
	where uuid = $1 and deleted is null;
`
