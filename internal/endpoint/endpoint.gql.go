package endpoint

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/gql"
	"switchboard/internal/scenario"
	"switchboard/internal/user"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var EndpointGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "EndpointInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"mockServiceId": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"path": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"method": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"responseDelay": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
	},
})

var EndpointGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Endpoint",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"path": &graphql.Field{
			Type: graphql.String,
		},
		"method": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"responseDelay": &graphql.Field{
			Type: graphql.Int,
		},
		"scenarios": &graphql.Field{
			Type: graphql.NewList(scenario.ScenarioGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				endpointID := p.Source.(Endpoint).ID
				scenarios, err := scenario.GetScenarios(endpointID)
				if err != nil {
					logrus.Errorln(err)
					return make([]scenario.Scenario, 0), gql.NewGqlError(common.ErrorGeneric, "could not retrieve scenarios")
				}
				return scenarios, nil
			},
		},
		"createdBy": &graphql.Field{
			Type: user.UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(Endpoint).CreatedBy
				users, err := user.GetUserByID(userId)
				if err != nil {
					logrus.Errorln(err)
					return make([]user.User, 0), gql.NewGqlError(common.ErrorGeneric, "could not resolve createdBy field")
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

func CreateEndpointResolver(p graphql.ResolveParams) (interface{}, error) {
	var input CreateEndpointRequestBody
	mapstructure.Decode(p.Args["endpoint"], &input)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	createdEndpoint, createErr := CreateEndpoint(currentUser.ID, &input)
	if createErr == nil {
		return *createdEndpoint, nil
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		return nil, gql.NewGqlError(common.ErrorDuplicateEntity, "duplicate endpoint")
	}

	return nil, gql.NewGqlError(common.ErrorGeneric, "could not create endpoint")
}

func DeleteEndpointResolver(p graphql.ResolveParams) (interface{}, error) {
	endpointID := p.Args["endpointId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	ok, err := DeleteEndpoint(currentUser.ID, endpointID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("could not delete endpoint %s. [error code: %s] [description: %s]", endpointID, err.ErrorCode, err.Description))
		return false, gql.NewGqlError(common.ErrorGeneric, "could not delete endpoint")
	}
	if !ok {
		return false, gql.NewGqlError(common.ErrorNotFound, "endpoint not found")
	}

	return true, nil
}
