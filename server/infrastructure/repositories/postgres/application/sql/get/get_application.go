package get

const ApplicationSQL = `
	with application_info_cte as (
		select uuid,
			   application_id,
			   package_name,
			   enterprise_id,
			   available,
			   location
		from public.applications
		where uuid = $1 and
			  enterprise_id = $2 and
			  deleted is null
	)
	select jsonb_build_object(
		'version_uuid', av.uuid,
		'application_uuid', aic.uuid,
		'package_name', aic.package_name,
		'enterprise_id', aic.enterprise_id,
		'available', aic.available,
		'location', aic.location,
		'metadata', jsonb_build_object(
			'uid', am.uuid
		),
		'google', jsonb_build_object(
			'uid', gai.uuid,
			'name', gai.name,
			'title', gai.title
		)
	)
	from public.application_versions av
		left join public.application_metadata am on am.metadata_id = av.application_metadata_id
		left join public.google_application_info gai on gai.id = av.google_application_info_id
		left join application_info_cte aic on aic.application_id = av.application_id
	where av.uuid = $3 and
		  av.application_id = (select application_id from application_info_cte) and
		  av.deleted is null;
`
