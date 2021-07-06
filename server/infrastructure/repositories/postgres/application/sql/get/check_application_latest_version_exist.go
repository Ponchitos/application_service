package get

const CheckApplicationLatestVersionExistSQL = `
	with max_version_cte as (
    	select av.version_id, max(av.version_code)
    	from public.applications a
             left join application_versions av on a.application_id = av.application_id
    	where a.uuid = $1 and
            a.enterprise_id = $2 and
        	a.deleted is null
    	group by av.version_id
    	limit 1
	)
	select version_id from max_version_cte;
`
