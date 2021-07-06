package get

const GoogleApplicationManagedPropertyByPropertyIDSQL = `
	select mp.property_id,
		   mp.key,
		   mp.type,
		   mp.title,
		   mp.description,
		   mp.default_value,
		   mp.entries
	from public.managed_properties_set mps
		left join public.managed_properties mp on mps.child_managed_property_id = mp.property_id and mp.deleted is null
	where mps.parent_managed_property_id = $1 and mps.deleted is null;
`
