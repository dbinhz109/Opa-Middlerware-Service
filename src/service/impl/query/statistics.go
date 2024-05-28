package query

import "go-app/src/service"

type Map = service.JsonMap

var APP_STATS_QUERY = Map{
	"version": true,
	"size":    0,
	"sort": []Map{
		{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"apps": Map{
			"terms": Map{
				"size":  100,
				"field": "detected_application_name.keyword",
			},
			"aggs": Map{
				"traffic_bytes": Map{
					"sum": Map{
						"field": "total_bytes",
					},
				},
				// "duration": Map{
				// 	"sum": Map{
				// 		"script": Map{"source": "doc['last_seen_at'].value - doc['first_seen_at'].value"},
				// 	},
				// },
				"bucket_sort": Map{
					"bucket_sort": Map{
						"sort": []Map{
							{"traffic_bytes": Map{"order": "desc"}},
						},
					},
				},
			},
		},
		"categories": Map{
			"terms": Map{
				"size":  100,
				"field": "category.application",
			},
			"aggs": Map{
				// "traffic_bytes": Map{
				// 	"sum": Map{
				// 		"field": "total_bytes",
				// 	},
				// },
				"duration": Map{
					"sum": Map{
						"script": Map{"source": "doc['last_seen_at'].value - doc['first_seen_at'].value"},
					},
				},
				"bucket_sort": Map{
					"bucket_sort": Map{
						"sort": []Map{
							{"duration": Map{"order": "desc"}},
						},
					},
				},
			},
		},
		"devices": Map{
			"terms": Map{
				"size":  100,
				"field": "local_mac.keyword",
			},
			"aggs": Map{
				"traffic_bytes": Map{
					"sum": Map{
						"field": "total_bytes",
					},
				},
				"bucket_sort": Map{
					"bucket_sort": Map{
						"sort": []Map{
							{"traffic_bytes": Map{"order": "desc"}},
						},
					},
				},
			},
		},
		"users": Map{
			"terms": Map{
				"size":  100,
				"field": "user_name.keyword",
			},
			"aggs": Map{
				"traffic_bytes": Map{
					"sum": Map{
						"field": "total_bytes",
					},
				},
				"bucket_sort": Map{
					"bucket_sort": Map{
						"sort": []Map{
							{"traffic_bytes": Map{"order": "desc"}},
						},
					},
				},
			},
		},
		"domains": Map{
			"terms": Map{
				"size":  100,
				"field": "dns_host_name.keyword",
			},
			"aggs": Map{
				"traffic_bytes": Map{
					"sum": Map{
						"field": "total_bytes",
					},
				},
				"bucket_sort": Map{
					"bucket_sort": Map{
						"sort": []Map{
							{"traffic_bytes": Map{"order": "desc"}},
						},
					},
				},
			},
		},
	},
	"query": Map{
		"bool": Map{
			"must": []Map{},
			"must_not": []Map{
				{"term": Map{"detected_protocol": 5}},
			},
		},
	},
}

var QUERY_ACTIVE_EDGES = Map{
	"size": 500,
	"sort": []Map{
		{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"query": Map{
		"bool": Map{
			"must": []Map{},
		},
	},
}

var QUERY_ACTIVE_DEVICES = Map{
	"size": 500,
	"sort": []Map{
		{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"grouping": Map{
			"multi_terms": Map{
				"terms": []Map{
					{"field": "edge_serial.keyword"},
					{"field": "local_mac.keyword"},
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

var CON_STATS_QUERY = Map{
	"version": true,
	"size":    0,
	"sort": []Map{
		{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{},
	"query": Map{
		"bool": Map{
			"must": []Map{},
		},
	},
}

var DURATION_STATS_QUERY = Map{
	"version": true,
	"size":    2,
	"sort": []Map{
		{
			"timestamp": Map{
				"order":         "asc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"apps": Map{
			"terms": Map{"field": "detected_application_name.keyword"},
			"aggs": Map{
				"domains": Map{
					"terms": Map{"field": "dns_host_name.keyword"},
					"aggs": Map{
						"duration": Map{
							"sum": Map{
								"script": Map{"source": "doc['last_seen_at'].value - doc['first_seen_at'].value"},
							},
						},
						"bucket_sort": Map{
							"bucket_sort": Map{
								"sort": []Map{
									{"duration": Map{"order": "desc"}},
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
			"must": []Map{
				{"range": Map{"timestamp": Map{"gte": "now-8h", "lte": "now"}}},
				{"term": Map{"other_type.keyword": Map{"value": "remote"}}},
			},
			"must_not": []Map{
				{"term": Map{"detected_protocol": 5}},
				{"term": Map{"detected_protocol": 81}},
			},
		},
	},
}

var QUERY_DAILY_APP_STATS = Map{
	"size":    1,
	"_source": false,
	"sort": []Map{
		{
			"timestamp": Map{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
	},
	"aggs": Map{
		"users": Map{
			"terms": Map{"field": "user_name.keyword", "size": 10000},
			"aggs": Map{
				"apps": Map{
					"terms": Map{"field": "app_name.keyword", "size": 10000},
					"aggs": Map{
						"duration": Map{
							"sum": Map{
								"field": "duration",
							},
						},
						"total_bytes": Map{
							"sum": Map{
								"field": "total_bytes",
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
