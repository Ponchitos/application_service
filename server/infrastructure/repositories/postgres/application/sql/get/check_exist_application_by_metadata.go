package get

const CheckExistApplicationByMetadataSQL = `
	select a.uuid, av.version_code
	from public.applications a
    	left join public.application_latest_versions alv on a.application_id = alv.application_id and alv.deleted is null
    	left join public.application_versions av on alv.version_id = av.version_id and av.deleted is null
	where a.package_name = (select package_name from application_metadata where metadata_id = $1 and deleted is null) and
      	  a.deleted is null;
`
