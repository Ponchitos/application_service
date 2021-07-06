package add

const ApplicationMetadataSQL = `
	insert into public.application_metadata
		(
			uuid, 
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
			uses_permissions
		)
	values
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16
		);
`
