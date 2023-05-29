func sendOKNotification(eq *ingvws.INGVEarthquake, w *_worker) error {
  maxdistance, _ := eqmath.PerceptibleUpTo(eq.Magnitude, eq.DepthKm,
  eqmath.PerceptibleMin)

  // calcola il bounding box del terremoto
  var eartquakeArea locationutils.BoundingBox = 
  locationutils.Squaring(eq.Origin, maxdistance)

  // inserisce nella tabella post_earthquake_ok_notification
  // i dispositivi nell'area attraverso la query vista sopra 
  err := RegisterUsersAndPhonesInDatabase(w, eartquakeArea, eq)
  if err != nil {
    return fmt.Errorf("can't register users and phones that are inside 
                      bounding box: %w", err)
  }
  devices,err :=w.scsdb.GetPhonesInEarthquakeBoundingBox(eartquakeArea)
  
  if err != nil {
    return fmt.Errorf("can't get phones inside bounding box: %w", err)
  }
  if len(devices) == 0 {
    w.log.Debug("no devices found in bounding box")
    return nil
  }
  // prepara il payload della notifica, servira' al client per 
  // comunicare al meglio con il server
  customDataMap := make(map[string]string)
  customDataMap := make(map[string]string)
  customDataMap["x-sc-category"] = "ok-notification"
  customDataMap["x-sc-earthquakeid"] = strconv.FormatInt(eq.ID, 10)
  customDataMap["x-sc-magnitude"] = 
    strconv.FormatFloat(eq.Magnitude, 'f', -1, 32)
  customDataMap["x-sc-latitude"] = 
    strconv.FormatFloat(float64(eq.Origin.Latitude), 'f', -1, 64)
  customDataMap["x-sc-longitude"] = 
    strconv.FormatFloat(float64(eq.Origin.Longitude), 'f', -1, 64)
  customDataMap["x-sc-name"] = eq.Location
  customDataMap["x-sc-location"] = eq.Location
  customDataMap["x-sc-intensity"] = 
    strconv.FormatFloat(eq.EstimatedIntensityAtEpicenter(), 'f', -1, 64)
  customDataMap["x-sc-range"] = 
    strconv.FormatFloat(eq.MaxDistance(), 'f', -1, 64)
  customDataMap["x-sc-depth"] = 
    strconv.FormatFloat(eq.DepthKm, 'f', -1, 64)
  // crea un nuovo oggetto di tipo Notification
  var okNotification = mobilepush.Notification{
      Title: "Ok notification", // inserisce il titolo della notifica
      Body:  "An high magnitude earthquake occurred in your position,
      are you ok?", // inserisce il body, ovvero una descrizione piu' 
      // dettagliata della notifica
      CustomData: customDataMap // inserisce il payload nella notifica
  }
  // divide i tokens raccolti dai telefoni per tipo "firebase" 
  // e "APNS" inserendoli all'interno dell'oggetto per notificarli
  okNotification.FirebaseTokens, okNotification.ApnsTokens = 
      GetFirebaseApnsTokens(devices)

  // manda la notifica push su tutti i token presenti nell'oggetto
  pushResults, err := w.push.SendPush(okNotification)
  if err != nil {
    return fmt.Errorf("can't send OK notification: %w", err)
  }

  // ora bisogna salvare nel database se sono avvenuti degli errori 
  // notificando i dispositivi. Pero' come e' stato detto in precedenza
  // di default nel database viene segnato errore all'inserimento della
  // tupla, in quanto la notifica non e' stata ancora inviata.
  // Infatti andremo a cambiare il valore a "0", ovvero nessun errore, 
  // solo per quei dispositivi che sono stati notificati correttamente
  for _, dev := range devices {
      var hasError bool
  
      switch dev.TypeToken {
      case "firebase":
          _, hasError = pushResults.FcmErrors[dev.NotificationToken]
      case "apns":
          _, hasError = pushResults.ApnsErrors[dev.NotificationToken]
      default:
          continue
      }

      if hasError {
          // if it has error no need to change because it is the 
          // database default value
          continue
      }
  // imposta a "no error" il valore nella tabella, per comunicare che
  // la notifica e' stata mandata correttamente
      err := w.scsdb.UpdateSendOkNotificationPushResult(dev.Phoneid, 
      eq.ID, hasError)
      if err != nil {
        w.log.WithError(err).Error("can't update the succeed of 
                                    push notifications")
      }
      
  }
}

// questa funzione ha l'unico scopo di incapsulare la query che aggiunge
// i telefoni colpiti da un terremoto nel database, dando cosi' un nome 
// piu' comprensibile nel contesto locale
func RegisterUsersAndPhonesInDatabase(w *_worker, box 
*locationutils.BoundingBox, quake *ingvws.INGVEarthquake) error {
    return w.scsdb.AddUsersPhonesToWatchingList(box, quake)
}