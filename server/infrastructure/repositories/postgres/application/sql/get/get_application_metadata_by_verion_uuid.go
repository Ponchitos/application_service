package get

const ApplicationMetadataByVersionUUIDSQL = `
	select am.uuid,
		   am.metadata_id,
		   am.link,
		   am.package_name,
		   am.application_label,
		   am.version_name,
		   am.file_size,
		   am.file_sha1_base64,
		   am.file_sha256_base64,
		   am.icon_base64,
		   am.externally_hosted_url,
		   am.native_codes,
		   am.certificate_base64s,
		   am.uses_features,
		   am.version_code,
		   am.maximum_sdk,
		   am.created,
		   am.uses_permissions
	from public.application_versions av
		left join public.application_metadata am on av.application_metadata_id = am.metadata_id and am.deleted is null
	where av.uuid = $1 and
		  av.application_id = (select application_id from public.applications where uuid = $2 and enterprise_id = $3 and deleted is null) and
		  av.deleted is null;
`
