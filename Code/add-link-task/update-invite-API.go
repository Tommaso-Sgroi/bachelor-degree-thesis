func (rt *_router) updateGroupInvitationID(w http.ResponseWriter, 
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

	invitationid, err := uuid.NewV4()
	if err != nil {
		// handle error
	}
	err = rt.DB.UpdateGroupInvitationID(group.ID, invitationid)
	if err != nil {
		// handle error
	}
	groupInvitation.InvitationID = invitationid.String()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(&groupInvitation)
}
