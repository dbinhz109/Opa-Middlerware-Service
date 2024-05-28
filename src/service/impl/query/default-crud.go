package query

var QUERY_CRUD_GETLIST = Map{
	"size": 5,
	"sort": []Map{
		{"timestamp": Map{
			"order":         "desc",
			"unmapped_type": "boolean",
		}},
	},
	"query": Map{
		"bool": Map{"must": []Map{}},
	},
}
