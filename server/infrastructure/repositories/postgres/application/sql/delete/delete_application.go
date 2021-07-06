package delete

const ApplicationSQL = `
	update public.applications
		set deleted = now()
	where uuid = $1 and
		  enterprise_id = $2
`
