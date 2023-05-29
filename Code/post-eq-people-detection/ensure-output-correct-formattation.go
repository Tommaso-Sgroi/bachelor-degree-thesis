type PeopleDetectionModelOutput struct {
	ModelName    string  `json:"modelName"`
	DetectedType string  `json:"detectedType"`
	Accuracy     float64 `json:"accuracy"`
}

func (p *PeopleDetectionModelOutput) Valid() bool {
  return p.Accuracy > 0.0 &&
         p.DetectedType != "" && 
         p.ModelName != ""
}