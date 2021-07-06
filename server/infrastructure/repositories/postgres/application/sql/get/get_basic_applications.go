package get

const BasicApplicationsSQL = `
	with applications_cte as (
		select a.uuid as "application_uuid",
			   av.created,
			   av.uuid as "version_uuid",
			   av.icon,
			   a.name,
			   a.package_name,
			   a.available,
			   a.location,
			   a.status,	
			   av.version_name,
			   av.version_code,
			   av.min_sdk
		from public.application_latest_versions alv
			left join public.applications a on a.application_id = alv.application_id and a.deleted is null
			left join public.application_versions av on av.version_id = alv.version_id and av.deleted is null
		where alv.enterprise_id = $1 and 
			  alv.deleted is null
	)
	select
		(select coalesce(count(application_uuid), 0) from applications_cte),
		jsonb_agg(jsonb_build_object(
			'uid', t.application_uuid,
			'created', t.created,
			'versionUid', t.version_uuid,
			'icon', t.icon,
			'name', t.name,
		    'status', t.status,
			'packageName', t.package_name,
			'available', t.available,
			'location', t.location,
			'versionName', t.version_name,
			'versionCode', t.version_code,
			'minSdk', t.min_sdk
		))
	from (select * from applications_cte offset $2 limit $3) as t;
	
`
