package gql

import (
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var UserGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func GetUserResolver(p graphql.ResolveParams) (interface{}, error) {
	userId, ok := p.Args["id"].(string)
	if ok {
		user, err := db.GetUserByID(userId)
		if err != nil {
			logrus.Errorln(err)
			return nil, NewGqlError(common.ErrorGeneric, "could not retrieve user")
		}
		return user, nil
	}
	return nil, nil
}

func GetUsersResolver(p graphql.ResolveParams) (interface{}, error) {
	users, err := db.GetUsers()
	if err != nil {
		logrus.Errorln(err)
		return make([]models.User, 0), NewGqlError(common.ErrorGeneric, "could not retrieve users")
	}
	return users, nil
}
