package generated

//go:generate oapi-codegen -generate types -package flights -o flights/types.go ../../../specs/openapi/flights_search.yaml
//go:generate oapi-codegen -generate types -package hotels -o hotels/types.go ../../../specs/openapi/hotel_search.yaml
//go:generate oapi-codegen -generate types -package hotellist -o hotellist/types.go ../../../specs/openapi/hotel_list.yaml
