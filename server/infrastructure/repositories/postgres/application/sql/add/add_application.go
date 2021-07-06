package add

const ApplicationSQL = `
	insert into public.applications
		(enterprise_id, package_name, available, location, name, status)
	values
		($1, $2, $3, $4, $5, $6)
	returning application_id, uuid;
`
