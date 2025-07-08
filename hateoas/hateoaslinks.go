package hateoas


import (
	"fmt"
	"myapi/models"
)

func CreateUserResponse(u models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":   u.ID,
        "name": u.Name,
        "age":  u.Age,
        "_links": map[string]string{
			"self":    fmt.Sprintf("/api/v1/users/%d", u.ID),
			"update": fmt.Sprintf("/api/v1/users/%d", u.ID),
			"delete": fmt.Sprintf("/api/v1/users/%d", u.ID),
            "allUsers": "/api/v1/users",
        },
    }
}
  