func (rt *_router) addPeopleDetectionResults(w http.ResponseWriter, 
    r *http.Request, p httprouter.Params, 
    ctx reqcontext.RequestContext) {
  if ctx.PhoneID == "" {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  
  var earthquakeID, err = strconv.ParseInt(p.ByName("earthquakeid"), 
                                           10, 64)
  if err != nil {
    ctx.Logger.WithError(err).Info("can't parse earthquake 
                                    id parameter")
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  exists, err := rt.DB.EarthquakeExists(earthquakeID)
  if err != nil {
    ctx.Logger.WithError(err).Debug("can't verify if an earthquake 
                                     id exists")
    w.WriteHeader(http.StatusBadRequest)
    return
  } else if !exists {
    ctx.Logger.WithError(err).Infof("earthquake id '%v' not found", 
                                    earthquakeID)
    w.WriteHeader(http.StatusNotFound)
    return
  }

  // decodifica il JSON 
  var peopleDetectionModelList = 
      make([]dbtypes.PeopleDetectionModelOutput, 6)
  err = json.NewDecoder(r.Body).Decode(&peopleDetectionModelList)
  if err != nil {
    ctx.Logger.WithError(err).Debug("cannot parse the request's body")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  // check if the json parameters are correct
  for _, p := range peopleDetectionModelList {
    if !p.Valid() {
        ctx.Logger.Debug("one or more fields are missing in people 
                          detection output json")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
  }
  
  // store output in database
  err = rt.DB.AddPeopleDetectionModelOutput(peopleDetectionModelList, 
                              ctx.PhoneID, ctx.UserID, earthquakeID)
  if err != nil {
    ctx.Logger.WithError(err).Info("can't add the model's output")
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusCreated)
}
