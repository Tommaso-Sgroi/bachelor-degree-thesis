func (rt *_router) getGroupInvitationID(w http.ResponseWriter, 
_ *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	if ctx.UserID == uuid.Nil && ctx.PhoneID == "" {
		// handle error
	}
	var group *database.Group
	var groupInvitation dbtypes.GroupInvitation
	var err error
	// User auth
	if ctx.UserID != uuid.Nil {
		group, err = rt.DB.GetMyGroup(ctx.UserID)
		if err != nil {
		// handle error
		}
	} else { // User not Auth
		group, err = rt.DB.GetPhoneGroup(ctx.PhoneID)
		if err != nil {
		// handle error
		}
	}

	if group.InvitationID == "" {
		// for retro-compatibility the invitation id can be void,
    // if it is we have to create a new one and insert it
		invitationid, err := uuid.NewV4()
		if err != nil {
		// handle error
		}
		err = rt.DB.UpdateGroupInvitationID(group.ID, invitationid)
		group.InvitationID = invitationid.String()
		if err != nil {
		// handle error
		}
	}

	if err != nil {
		// handle error
	}
	// transfer data from group's struct in GroupInvitation's struct
	groupInvitation.InvitationID = group.InvitationID

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&groupInvitation)
}
