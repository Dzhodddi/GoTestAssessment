package graphql

import (
	"github.com/graphql-go/graphql"
)

type SpyCatInfo struct {
	ID               int64
	Name             string
	YearOfExperience int
	Breed            string
	Salary           int
	Mission          Mission
}

type Mission struct {
	ID       int64
	Complete bool
	Targets  []Target
}

type Target struct {
	ID       int64
	Name     string
	Country  string
	Notes    string
	Complete bool
}

var targetType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Target",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
		},
		"Name": &graphql.Field{
			Type: graphql.String,
		},
		"Country": &graphql.Field{
			Type: graphql.String,
		},
		"Notes": &graphql.Field{
			Type: graphql.String,
		},
		"Complete": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var missionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mission",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
		},
		"Complete": &graphql.Field{
			Type: graphql.Boolean,
		},
		"Targets": &graphql.Field{
			Type: graphql.NewList(targetType),
		},
	},
})

var spyCatInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SpyCatInfo",
	Fields: graphql.Fields{
		"ID": &graphql.Field{
			Type: graphql.Int,
		},
		"Name": &graphql.Field{
			Type: graphql.String,
		},
		"YearOfExperience": &graphql.Field{
			Type: graphql.Int,
		},
		"Breed": &graphql.Field{
			Type: graphql.String,
		},
		"Salary": &graphql.Field{
			Type: graphql.Int,
		},
		"Mission": &graphql.Field{
			Type: missionType,
		},
	},
})

type Cat struct {
	resolver *Resolver
}

func (schema Cat) NewCatSchema() graphql.Schema {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"list": &graphql.Field{
				Type:    graphql.NewList(spyCatInfoType),
				Resolve: schema.resolver.getOneCat,
			},
			"cat": &graphql.Field{
				Type: spyCatInfoType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: schema.resolver.getOneCat,
			},
		},
	})
	catSchema, _ := graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})
	return catSchema
}

func (schema Cat) GetListOfCats() *graphql.Result {
	query := `{
		  list {
			ID
			Name
			YearOfExperience
			Breed
			Salary
			Mission {
			  ID
			  Complete
			  Targets {
				ID
				Name
				Country
				Notes
				Complete
			  }
			}
		  }
		}
	`
	data := graphql.Do(graphql.Params{
		Schema:        schema.NewCatSchema(),
		RequestString: query,
	})
	return data
}
