package model

type CreateKosRequest struct {
	Name string `json:"name" validate:required`
	OwnerID string `json:"owner_id" validate:required`
	
	
}

// type CreateKosResponse struct {
	
// }

// type Kos struct {
// 		OwnerID string
//         Name	string
//         Description
//         StreetAddress
//         DistrictID
//         PostalCode
//         Latitude
//         Longitude
//         GenderType
//         Status
//         IsVerified
//         CreatedAt
//         UpdatedAt
// }