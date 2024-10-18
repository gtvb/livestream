package http

import (
	"github.com/gtvb/livestream/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ###### USER RELATED TYPES ######

type FollowBody struct {
	// User's id (this is the one that is performing the follow operation)
	UserID primitive.ObjectID
}

// FollowUserParamsWrapper contains parameters for updating a user
// swagger:parameters followUser
type FollowUserParamsWrapper struct {
	// in:body
	Body FollowBody
}

type UpdateUserBody struct {
	// User's username
	// required: true
	Username string `json:"username"`
	// User's email
	// required: true
	Email string `json:"email"`
	// User's password
	// required: true
	Password string `json:"password"`
}

// UpdateUserParamsWrapper contains parameters for updating a user
// swagger:parameters updateUser
type UpdateUserParamsWrapper struct {
	// in:body
	Body UpdateUserBody
}

type LoginBody struct {
	// User's email
	// required: true
	Email string `json:"email"`
	// User's password
	// required: true
	Password string `json:"password"`
}

// LoginParamsWrapper contains parameters for user login.
// swagger:parameters loginUser
type LoginParamsWrapper struct {
	// in:body
	Body LoginBody
}

type SignupBody struct {
	// User's name
	// required: true
	Name string `json:"name"`
	// User's username
	// required: true
	Username string `json:"username"`
	// User's email
	// required: true
	Email string `json:"email"`
	// User's password
	// required: true
	Password string `json:"password"`
}

// SignupParamsWrapper contains parameters for user signup.
// swagger:parameters signupUser
type SignupParamsWrapper struct {
	// in:body
	Body SignupBody
}

// UserResponseWrapper contains a user response.
// swagger:response userResponse
type UserResponseWrapper struct {
	// in:body
	Body struct {
		// The user details
		User models.User `json:"user"`
	}
}

// UserListResponseWrapper contains a user list response.
// swagger:response userListResponse
type UserListResponseWrapper struct {
	// in:body
	Body struct {
		// The user details
		Users []models.User `json:"users"`
	}
}

// LiveStreamsResponseWrapper contains a response with live streams.
// swagger:response liveStreamsResponse
type LiveStreamsResponseWrapper struct {
	// in:body
	Body struct {
		// List of live streams
		LiveStreams []models.LiveStream `json:"livestreams"`
	}
}

// TokenResponseWrapper contains a token response.
// swagger:response tokenResponse
type TokenResponseWrapper struct {
	// The JWT token for future protected requests.
	// required: true
	Body struct {
		Token string `json:"token"`
	}
}

// MessageResponseWrapper contains a message response.
// swagger:response messageResponse
type MessageResponseWrapper struct {
	Body struct {
		// A descriptive message
		Message string `json:"message"`
	}
}

// ###### LIVESTREAM RELATED TYPES ######

type UpdateLiveStreamBody struct {
	// Live Status. On or off
	// required: true
	LiveStatus *bool `json:"live_stream_status"`
	// Name of the live stream
	// required: true
	Name string `json:"name"`
}

// UpdateLiveStreamParamsWrapper contains parameters for updating a live stream.
// swagger:parameters updateLiveStream
type UpdateLiveStreamParamsWrapper struct {
	// in:body
	Body UpdateLiveStreamBody
}

type CreateLiveStreamBody struct {
	// User ID of the stream creator
	// required: true
	UserId string `json:"user_id"`
	// Name of the live stream
	// required: true
	Name string `json:"name"`
	// User password (unhashed, obviously)
	// required: true
	Password string `json:"password"`
}

// CreateLiveStreamParamsWrapper contains parameters for creating a live stream.
// swagger:parameters createLiveStream
type CreateLiveStreamParamsWrapper struct {
	// in:body
	Body CreateLiveStreamBody
}

// LiveStreamResponseWrapper contains a response with live stream data.
// swagger:response liveStreamResponse
type LiveStreamResponseWrapper struct {
	// in:body
	Body struct {
		// ID of the live stream
		StreamId primitive.ObjectID `json:"stream_id"`
	}
}
