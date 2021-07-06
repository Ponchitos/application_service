package edit

const RollbackUpdateApplicationStatusByVersionUUIDSQL = `
	update public.applications
		set status = previous_status,
		    modified = now()
	where application_id = (select application_id from public.application_versions where uuid = $1 and deleted is null) and
		  enterprise_id = $2 and
		  deleted is null;
`
