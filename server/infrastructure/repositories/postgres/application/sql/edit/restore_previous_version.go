package edit

const RestorePreviousVersionSQL = `
	with previous_version_cte as (
		select set_id
		from public.application_group_set
		where group_uuid = $1 and 
		  application_version_id != (select version_id from public.application_versions where uuid = $2 and deleted is null) and
		  application_id = (select application_id from public.applications where uuid = $3 and deleted is null) and
		  enterprise_id = $4 and
		  deleted is not null
		order by deleted desc
		limit 1
	)
	update previous_version_cte
		set deleted = null
	where set_id = (select set_id previous_version_cte);
`
