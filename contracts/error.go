package contracts

type HttpError struct {
	Message    string `json:"message"`
	IncidentId string `json:"incidentId,omitempty"`
}
