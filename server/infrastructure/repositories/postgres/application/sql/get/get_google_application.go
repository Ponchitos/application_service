package get

const GoogleApplicationSQL = `
	select id,
		   uuid,
		   name,
		   title	
	from public.google_application_info gai
	where gai.uuid = $1;
`
