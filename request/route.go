package request

type CreateRouteInput struct {
	UserID              uint   `json:"userID"`
	OriginCityName      string `json:"originCityName"`
	DestinationCityName string `json:"destinationCityName"`
	Destinations        []uint `json:"destinations"`
}
