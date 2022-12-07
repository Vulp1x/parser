package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/pb/instaproxy"
)

type InstUser struct {
	Pk                       int64       `json:"pk"`
	Username                 string      `json:"username"`
	FullName                 string      `json:"full_name"`
	IsPrivate                bool        `json:"is_private"`
	ProfilePicUrl            string      `json:"profile_pic_url"`
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

type InstUserShort struct {
	Pk              int64         `json:"pk"`
	Username        string        `json:"username"`
	FullName        string        `json:"full_name"`
	ProfilePicUrl   string        `json:"profile_pic_url"`
	ProfilePicUrlHd string        `json:"profile_pic_url_hd"`
	IsPrivate       bool          `json:"is_private"`
	IsVerified      bool          `json:"is_verified"`
	Stories         []interface{} `json:"stories"`
}

func FullUserFromProto(protoUser *instaproxy.FullUser) *InstUser {
	u := &InstUser{}
	if protoUser == nil {
		return u
	}

	u.Pk = protoUser.Pk
	u.Username = protoUser.Username
	u.FullName = protoUser.FullName
	u.IsPrivate = protoUser.IsPrivate
	u.ProfilePicUrl = protoUser.ProfilePicUrl
	u.IsVerified = protoUser.IsVerified
	u.MediaCount = int(protoUser.MediaCount)
	u.FollowerCount = int(protoUser.FollowerCount)
	u.FollowingCount = int(protoUser.FollowingCount)
	u.Biography = protoUser.Biography
	u.ExternalUrl = &protoUser.ExternalUrl
	u.IsBusiness = protoUser.IsBusiness
	// u.PublicEmail  =
	// u.ContactPhoneNumber =
	// u.PublicPhoneCountryCode  =
	// u.PublicPhoneNumber       =

	return u
}

func ShortUsersFromProto(users []*instaproxy.UserShort) ShortInstUsers {
	var domainUsers = make([]InstUserShort, 0, len(users))
	for _, userShort := range users {
		if userShort == nil {
			continue
		}

		domainUsers = append(domainUsers, InstUserShort{
			Pk:            int64(userShort.Pk),
			Username:      userShort.Username,
			FullName:      userShort.FullName,
			ProfilePicUrl: userShort.ProfilePicUrl,
			IsPrivate:     userShort.IsPrivate,
			IsVerified:    userShort.IsVerified,
		})
	}

	return domainUsers
}

func (u InstUser) ToUpdateParams(id uuid.UUID, isCorrect bool) dbmodel.UpdateBloggerParams {
	parsedAt := time.Now()
	return dbmodel.UpdateBloggerParams{
		UserID:                 u.Pk,
		FollowersCount:         int64(u.FollowerCount),
		ParsedAt:               &parsedAt,
		IsCorrect:              isCorrect,
		IsPrivate:              u.IsPrivate,
		IsVerified:             u.IsVerified,
		IsBusiness:             u.IsBusiness,
		FollowingsCount:        int32(u.FollowingCount),
		ContactPhoneNumber:     u.ContactPhoneNumber,
		PublicPhoneNumber:      u.PublicPhoneNumber,
		PublicPhoneCountryCode: u.PublicPhoneCountryCode,
		CityName:               u.CityName,
		PublicEmail:            u.PublicEmail,
		ID:                     id,
	}
}

type InstUsers []InstUser

func (u InstUsers) ToSaveBloggersParmas(datasetID uuid.UUID) []dbmodel.SaveBloggersParams {
	params := make([]dbmodel.SaveBloggersParams, len(u))

	for i, user := range u {
		params[i] = dbmodel.SaveBloggersParams{
			DatasetID:              datasetID,
			Username:               user.Username,
			UserID:                 user.Pk,
			FollowersCount:         int64(user.FollowerCount),
			FollowingsCount:        int32(user.FollowingCount),
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

func (u InstUserShort) ToUpdateParams(id uuid.UUID, isCorrect bool) dbmodel.UpdateBloggerParams {
	parsedAt := time.Now()
	return dbmodel.UpdateBloggerParams{
		UserID:     u.Pk,
		ParsedAt:   &parsedAt,
		IsCorrect:  isCorrect,
		IsPrivate:  u.IsPrivate,
		IsVerified: u.IsVerified,
		ID:         id,
	}
}

type ShortInstUsers []InstUserShort

func (su ShortInstUsers) ToSaveBloggersParmas(datasetID uuid.UUID) []dbmodel.SaveBloggersParams {
	params := make([]dbmodel.SaveBloggersParams, len(su))
	for i, user := range su {
		params[i] = dbmodel.SaveBloggersParams{
			DatasetID:  datasetID,
			Username:   user.Username,
			UserID:     user.Pk,
			IsPrivate:  user.IsPrivate,
			IsVerified: user.IsVerified,
		}
	}

	return params
}

func (su ShortInstUsers) ToSaveTargetsParams(mediaPk int64, datasetID uuid.UUID) dbmodel.SaveTargetUsersParams {
	var usernames = make([]string, len(su))
	var fullNames = make([]string, len(su))
	var userIDs = make([]int64, len(su))
	var isPrivates = make([]bool, len(su))
	var isVerified = make([]bool, len(su))

	for i, user := range su {
		usernames[i] = user.Username
		userIDs[i] = user.Pk
		fullNames[i] = user.FullName
		isPrivates[i] = user.IsPrivate
		isVerified[i] = user.IsVerified
	}

	return dbmodel.SaveTargetUsersParams{
		Usernames:  usernames,
		UserIds:    userIDs,
		FullNames:  fullNames,
		IsPrivate:  isPrivates,
		IsVerified: isVerified,
		MediaPk:    mediaPk,
		DatasetID:  datasetID,
	}
}
