package user

import (
	"time"

	"go.uber.org/fx"
)

// Module provides all constructor and invocation methods to facilitate credits module
var Module = fx.Options(
	fx.Provide(
		NewDBRepository,
		NewService,
	),
)

type (
	User struct {
		tableName      struct{}    `pg:"user,discard_unknown_columns"`
		ID             int         `json:"id" pg:"id"`
		FirstName      string      `json:"first_name" pg:"first_name"`
		LastName       string      `json:"last_name" pg:"last_name"`
		Mobile         string      `json:"mobile" pg:"mobile,unique"`
		ProfilePicture string      `json:"profile_picture" pg:"profile_picture"`
		DOB            *time.Time  `json:"dob" form:"dob" time_format:"2006-01-02" pg:"dob"`
		CreatedAt      *time.Time  `json:"created_at" form:"created_at" pg:"created_at"`
		UpdatedAt      *time.Time  `json:"updated_at" form:"updated_at" pg:"updated_at"`
		Metadata       interface{} `json:"metadata,omitempty" pg:"metadata,type:jsonb"`
	}

	Pagination struct {
		CurrentPage    int `json:"current_page,omitempty"`
		TotalPages     int `json:"total_pages,omitempty"`
		TotalDataCount int `json:"total_data_count,omitempty"`
	}

	UserRequest struct {
		Mobile *string `form:"mobile,omitempty"`
		Name   *string `form:"name,omitempty"`
		Page   int     `form:"page,default=1"`
		Limit  int     `form:"limit,default=20"`
	}
)
