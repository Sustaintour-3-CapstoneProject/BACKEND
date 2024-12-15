package request

type CreateRouteInput struct {
	UserID              uint    `json:"userID"`
	OriginCityName      string  `json:"originCityName"`
	DestinationCityName string  `json:"destinationCityName"`
	Destinations        []uint  `json:"destinations"`
	Distance            float64 `json:"distance"`
	Time                string  `json:"time"`
	Cost                int     `json:"cost"`
}
