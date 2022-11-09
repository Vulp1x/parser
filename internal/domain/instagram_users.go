package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
)

type InstUser struct {
	Pk                       int64       `json:"pk"`
	Username                 string      `json:"username"`
	FullName                 string      `json:"full_name"`
	IsPrivate                bool        `json:"is_private"`
	ProfilePicUrl            string      `json:"profile_pic_url"`
	ProfilePicUrlHd          string      `json:"profile_pic_url_hd"`
	IsVerified               bool        `json:"is_verified"`
	MediaCount               int         `json:"media_count"`
	FollowerCount            int         `json:"follower_count"`
	FollowingCount           int         `json:"following_count"`
	Biography                string      `json:"biography"`
	ExternalUrl              *string     `json:"external_url"`
	AccountType              *int        `json:"account_type"`
	IsBusiness               bool        `json:"is_business"`
	PublicEmail              *string     `json:"public_email"`
	ContactPhoneNumber       *string     `json:"contact_phone_number"`
	PublicPhoneCountryCode   *string     `json:"public_phone_country_code"`
	PublicPhoneNumber        *string     `json:"public_phone_number"`
	BusinessContactMethod    string      `json:"business_contact_method"`
	BusinessCategoryName     string      `json:"business_category_name"`
	CategoryName             interface{} `json:"category_name"`
	Category                 interface{} `json:"category"`
	AddressStreet            interface{} `json:"address_street"`
	CityId                   interface{} `json:"city_id"`
	CityName                 *string     `json:"city_name"`
	Latitude                 interface{} `json:"latitude"`
	Longitude                interface{} `json:"longitude"`
	Zip                      string      `json:"zip"`
	InstagramLocationId      interface{} `json:"instagram_location_id"`
	InteropMessagingUserFbid interface{} `json:"interop_messaging_user_fbid"`
}

type InstUsers []InstUser

func (u InstUsers) ToSaveBloggersParmas(datasetID uuid.UUID) []dbmodel.SaveBloggersParams {
	params := make([]dbmodel.SaveBloggersParams, len(u))
	parsedAt := time.Now()

	for i, user := range u {
		params[i] = dbmodel.SaveBloggersParams{
			DatasetID:              datasetID,
			Username:               user.Username,
			UserID:                 user.Pk,
			FollowersCount:         int64(user.FollowerCount),
			FollowingsCount:        int32(user.FollowingCount),
			IsInitial:              false,
			ParsedAt:               &parsedAt,
			Parsed:                 true,
			IsPrivate:              user.IsPrivate,
			IsVerified:             user.IsVerified,
			IsBusiness:             user.IsBusiness,
			ContactPhoneNumber:     user.ContactPhoneNumber,
			PublicPhoneNumber:      user.PublicPhoneNumber,
			PublicPhoneCountryCode: user.PublicPhoneCountryCode,
			CityName:               user.CityName,
			PublicEmail:            user.PublicEmail,
		}
	}

	return params
}
