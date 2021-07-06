package get

const GoogleApplicationInfoByVersionUUIDSQL = `
	with google_version_info_cte as (
		select google_application_info_id
		from public.application_versions
		where uuid = $1 and
			  application_id = (
				select application_id 
				from public.applications 
				where uuid = $2 and 
					  enterprise_id = $3 and 
					  deleted is null
			  ) and
		      deleted is null
	)
	select jsonb_build_object(
		'uid', uuid,
		'name', name,
		'title', title,
		'status', status
	)
	from public.google_application_info
	where id = (select google_application_info_id from google_version_info_cte)	and
		  deleted is null;
`
