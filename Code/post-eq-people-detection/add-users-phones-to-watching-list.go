func (db *appdbimpl) AddUsersPhonesToWatchingList(
  box *locationutils.BoundingBox, quake *ingvws.INGVEarthquake
) error {
  _, err := db.c.Exec(`
INSERT INTO post_earthquake_ok_notification (userid, deviceid, 
    earthquakeid, eventended, owner)
SELECT DISTINCT d.userid, d.deviceid, ?, 
    IF(EXISTS(
        SELECT pt.phoneid FROM phones_tokens as pt 
        WHERE pt.phoneid = d.owner
        ), 0, 1), d.owner 
-- inserisce 0 per comunicare che l'evento NON e' chiuso
-- inserisce 1 per comunicare che l'evento e' chiuso
FROM devices d
WHERE d.latitude >= ?/*low_lat*/ AND d.latitude <= ?/*max_lat*/
  AND d.longitude >= ?/*low_long*/ AND d.longitude <= ?/*max_long*/;
`,  quake.ID, box.A.Latitude, box.B.Latitude, box.A.Longitude, 
    box.B.Longitude) // valori passati ai placeholder della query
  return err
}