package query

var AllAlerts = Map{
	"version": true,
	"size":    500,
	"sort": []any{
		Map{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"hourly": Map{
			"date_histogram": Map{
				"field":          "timestamp",
				"fixed_interval": "1h",
				"time_zone":      "Asia/Ho_Chi_Minh",
				"min_doc_count":  1,
			},
		},
		"groups": Map{
			"terms": Map{
				"field":         "rule.groups",
				"min_doc_count": 1,
			},
		},
		"level": Map{
			"range": Map{
				"field": "rule.level",
				"ranges": []Map{
					{"to": 11},
					{"from": 12},
				},
			},
		},
	},
	"_source": Map{
		"excludes": []any{
			"@timestamp",
		},
	},
	"query": Map{
		"bool": Map{
			"must": []any{
				Map{
					"range": Map{
						"timestamp": Map{},
					},
				},
			},
			"must_not": []any{},
			"should":   []any{},
		},
	},
}
var FindActiveAgents = Map{
	"version": true,
	"size":    1,
	"sort": []any{
		Map{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"_source": Map{
		"excludes": []any{
			"@timestamp",
		},
	},
	"aggs": Map{
		"agents": Map{
			"terms": Map{
				"field": "agent.id",
				"size":  200,
			},
		},
	},
	"query": Map{
		"bool": Map{
			"must": []any{},
			// "must_not": []any{},
			// "should":   []any{},
		},
	},
}

var EventAlerts = Map{
	"version": true,
	"size":    500,
	"sort": []any{
		Map{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"_source": Map{
		"excludes": []any{
			"@timestamp",
		},
	},
	"query": Map{
		"bool": Map{
			"must": []Map{},
		},
	},
}

var EventIbfd = Map{
	"version": true,
	"size":    500,
	"sort": []any{
		Map{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"key": Map{
			"multi_terms": Map{
				"terms": []Map{
					{
						"field": "local_peer.keyword",
					},
					{
						"field": "remote_peer.keyword",
					},
					{
						"field": "name.keyword",
					},
				},
			},
			"aggs": Map{
				"data": Map{
					"top_hits": Map{
						"size": 1,
						"sort": []Map{
							{
								"timestamp": Map{
									"order": "desc",
								},
							},
						},
					},
				},
			},
		},
	},
	"query": Map{
		"bool": Map{
			"must": []Map{},
		},
	},
}

var EventIbfdAvg = Map{
	"version": true,
	"size":    500,
	"sort": []any{
		Map{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"key": Map{
			"multi_terms": Map{
				"terms": []Map{
					{
						"field": "local_peer.keyword",
					},
					{
						"field": "remote_peer.keyword",
					},
					{
						"field": "name.keyword",
					},
				},
			},
			"aggs": Map{
				"data": Map{
					"date_histogram": Map{
						"field":             "timestamp",
						"calendar_interval": "1h",
					},
					"aggs": Map{
						"loss": Map{
							"avg": Map{
								"field": "loss",
							},
						},
						"jitter": Map{
							"avg": Map{
								"field": "jitter",
							},
						},
						"latency": Map{
							"avg": Map{
								"field": "latency",
							},
						},
					},
				},
			},
		},
	},
	"query": Map{
		"bool": Map{
			"must": []Map{},
		},
	},
}
