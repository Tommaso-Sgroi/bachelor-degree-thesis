func (rt *_router) redeemGroupInvite(w http.ResponseWriter, 
rq *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	if ctx.UserID == uuid.Nil && ctx.PhoneID == "" {
		// handle error
	}
	var invitationCode dbtypes.GroupInvitationRedeem
	err := json.NewDecoder(rq.Body).Decode(&invitationCode)
	if err != nil {
		// handle error
	}

	groupid, err := rt.DB.GetGroupFromInvitationID(&invitationCode)
	if err != nil {
		// handle error
	}

	// Check if the user or the phones aren't already in another group
	if ctx.UserID != uuid.Nil {
		group, err := rt.DB.GetMyGroup(ctx.UserID)
		if group != nil {
		// handle error
		}
		success, err := rt.DB.AddUserToGroup(ctx.UserID, groupid)
		switch {
		case err != nil:
		// handle error
		case !success:
		// handle error
		default:
			w.WriteHeader(http.StatusOK)
		}
	} else if ctx.PhoneID != "" {
		group, err := rt.DB.GetPhoneGroup(ctx.PhoneID)
		if group != nil {
		// handle error
		}
		success, err := rt.DB.AddPhoneToGroup(ctx.PhoneID, groupid)
		switch {
		case err != nil:
		// handle error
		case !success:
		// handle error
		default:
			w.WriteHeader(http.StatusOK)
		}
	}
}

