package get

const GoogleApplicationTracksSQL = `
	select track_id, track_alias
	from public.application_tracks
	where google_application_info_id = $1 and
		  deleted is null;	

`
