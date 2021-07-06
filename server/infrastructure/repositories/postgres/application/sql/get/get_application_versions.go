package get

const ApplicationVersionsSQL = `
	select jsonb_agg(jsonb_build_object(
		'uid', a.uuid,
		'created', av.created,
		'versionUid', av.uuid,
		'icon', av.icon,
		'name', a.name,
		'packageName', a.package_name,
		'available', a.available,
		'location', a.location,
		'versionName', av.version_name,
		'versionCode', av.version_code,
		'minSdk', av.min_sdk
	))
	from public.applications a
		left join public.application_versions av on av.application_id = a.application_id and av.deleted is null
	where a.uuid = $1 and
		  a.enterprise_id = $2 and
		  a.deleted is null;
`
