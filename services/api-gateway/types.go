package main

import "ride-sharing/shared/types"

type PreviewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type GetRouteRequest struct {
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type Driver struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlate"`
	PackageSlug    string `json:"packageSlug"`
}
