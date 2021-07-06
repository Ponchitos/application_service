package get

const CheckDeleteApplicationAvailableSQL = `
	select a.location,
		   a.status
	from public.application_versions av
		left join public.applications a on a.application_id = av.application_id and a.deleted is null
	where av.uuid = $1 and
		  a.uuid = $2 and
		  a.enterprise_id = $3 and
		  av.deleted is null;
`
