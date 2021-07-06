package get

const GoogleApplicationManagedPropertyByAppIDSQL = `
	select property_id,
		   key,
		   type,
		   title,
		   description,
		   default_value,
		   entries
	from public.managed_properties
	where google_application_info_id = $1 and deleted is null;
`
