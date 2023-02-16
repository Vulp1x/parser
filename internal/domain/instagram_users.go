package domain

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/pb/instaproxy"
)

type BioLink struct {
	LinkId                          int64  `json:"link_id"`
	Url                             string `json:"url"`
	LynxUrl                         string `json:"lynx_url"`
	LinkType                        string `json:"link_type"`
	Title                           string `json:"title"`
	GroupId                         int    `json:"group_id"`
	OpenExternalUrlWithInAppBrowser bool   `json:"open_external_url_with_in_app_browser"`
}

type FullUser dbmodel.FullTarget

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

func FullUserFromProto(protoUser *instaproxy.FullUser) *FullUser {
	if protoUser == nil {
		return &FullUser{}
	}

	user := &FullUser{
		ID:                         uuid.UUID{},
		DatasetID:                  uuid.UUID{},
		ParsedAt:                   time.Time{},
		Username:                   protoUser.Username,
		InstPk:                     protoUser.Pk,
		FullName:                   protoUser.FullName,
		IsPrivate:                  protoUser.IsPrivate,
		IsVerified:                 protoUser.IsVerified,
		IsBusiness:                 protoUser.IsBusiness,
		IsPotentialBusiness:        protoUser.IsPotentialBusiness,
		HasAnonymousProfilePicture: protoUser.HasAnonymousProfilePicture,
		Biography:                  protoUser.Biography,
		ExternalUrl:                protoUser.ExternalUrl,
		MediaCount:                 protoUser.MediaCount,
		FollowerCount:              protoUser.FollowerCount,
		FollowingCount:             protoUser.FollowingCount,
		Category:                   protoUser.Category,
		CityName:                   protoUser.CityName,
		ContactPhoneNumber:         protoUser.ContactPhoneNumber,
		Latitude:                   protoUser.Latitude,
		Longitude:                  protoUser.Longitude,
		PublicEmail:                protoUser.PublicEmail,
		PublicPhoneCountryCode:     protoUser.PublicPhoneCountryCode,
		PublicPhoneNumber:          protoUser.PublicPhoneNumber,
		WhatsappNumber:             protoUser.WhatsappNumber,
	}

	if len(protoUser.BioLinks) != 0 {
		bioLonksBytes, _ := json.Marshal(protoUser.BioLinks)
		user.BioLinks = string(bioLonksBytes)
	}

	return user
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

func (u FullUser) ToUpdateParams(id uuid.UUID) dbmodel.UpdateBloggerParams {
	parsedAt := time.Now()
	return dbmodel.UpdateBloggerParams{
		UserID:    u.InstPk,
		ParsedAt:  &parsedAt,
		IsCorrect: true,
		ID:        id,
	}
}

type FullUsers []FullUser

func (u FullUser) ToSaveFullTargetParams(datasetID uuid.UUID) dbmodel.SaveFullTargetParams {

	return dbmodel.SaveFullTargetParams{
		DatasetID:                  datasetID,
		Username:                   u.Username,
		InstPk:                     u.InstPk,
		FullName:                   u.FullName,
		IsPrivate:                  u.IsPrivate,
		IsVerified:                 u.IsVerified,
		IsBusiness:                 u.IsBusiness,
		IsPotentialBusiness:        u.IsPotentialBusiness,
		HasAnonymousProfilePicture: u.HasAnonymousProfilePicture,
		Biography:                  u.Biography,
		ExternalUrl:                u.ExternalUrl,
		MediaCount:                 u.MediaCount,
		FollowerCount:              u.FollowerCount,
		FollowingCount:             u.FollowingCount,
		Category:                   u.Category,
		CityName:                   u.CityName,
		ContactPhoneNumber:         u.ContactPhoneNumber,
		Latitude:                   u.Latitude,
		Longitude:                  u.Longitude,
		PublicEmail:                u.PublicEmail,
		PublicPhoneCountryCode:     u.PublicPhoneCountryCode,
		PublicPhoneNumber:          u.PublicPhoneNumber,
		BioLinks:                   u.BioLinks,
		WhatsappNumber:             u.WhatsappNumber,
	}
}

func (u FullUser) Format(format int) string {
	b := bytes.Buffer{}
	switch format {
	case 1:
		b.WriteString(strconv.FormatInt(u.InstPk, 10))
	case 2:
		b.WriteString(u.Username)
	case 3:

		b.WriteString(u.Username)
		b.WriteByte(',')
		b.WriteString(strconv.FormatInt(u.InstPk, 10))
	default:
		return ""
	}

	b.WriteByte(',')
	b.WriteString(u.ContactPhoneNumber)
	return b.String()
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

// Ptr возвращает указать на переданное значение
func Ptr[T any](val T) *T { return &val }
