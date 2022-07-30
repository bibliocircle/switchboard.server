package gql

import (
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var WorkspaceGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Workspace",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"mockServices": &graphql.Field{
			Type:    graphql.NewList(MockServiceGqlType),
			Resolve: GetWorkspaceMockServicesResolver,
		},
		"expiresAt": &graphql.Field{
			Type: graphql.String,
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(models.Workspace).CreatedBy
				users, err := db.GetUserByID(userId)
				if err != nil {
					logrus.Errorln(err)
					return make([]models.User, 0), NewGqlError(GqlInternalError, "could not resolve createdBy field")
				}
				return users, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func GetWorkspacesResolver(p graphql.ResolveParams) (interface{}, error) {
	wss, err := db.GetWorkspaces()
	if err != nil {
		logrus.Errorln(err)
		return make([]models.Workspace, 0), NewGqlError(GqlInternalError, "could not retrieve workspaces")
	}
	return wss, nil
}

func GetUserWorkspacesResolver(p graphql.ResolveParams) (interface{}, error) {
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	wss, err := db.GetUserWorkspaces(currentUser.ID)
	if err != nil {
		logrus.Errorln(err)
		return make([]models.Workspace, 0), NewGqlError(GqlInternalError, "could not retrieve user workspaces")
	}
	return wss, nil
}

func GetUserWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID, ok := p.Args["workspaceId"].(string)
	if ok {
		currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
		wss, err := db.GetUserWorkspaceByID(currentUser.ID, workspaceID)
		if err != nil {
			return nil, err
		}
		return *wss, nil
	}
	return nil, nil
}
