package edit

const SetUpdateAvailableStatusToApplicationSQL = `
	update public.applications
		set status = 'UPDATE_AVAILABLE'
	where uuid = $1 and
		  enterprise_id = $2 and
		  status = 'INSTALLED';
`
