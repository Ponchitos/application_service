package get

const GoogleAppPermissionsSQL = `
	select name,
		   description,
		   permission_id	
	from public.application_permissions
	where google_application_info_id = $1 and deleted is null;
`
