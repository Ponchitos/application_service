package get

const LatestApplicationNameByPackageNameSQL = `
	select
		av.version_code
	from public.application_latest_versions alv
		left join public.application_versions av on av.version_id = alv.version_id and av.deleted is null 
	where alv.application_id = (select application_id from public.applications where package_name = $1 and deleted is null) and
		  alv.deleted is null;
`
