package query

var QUERY_APP_USAGE = Map{
	"version": true,
	"size":    1,
	"sort": []Map{
		{"timestamp": Map{"order": "asc", "unmapped_type": "boolean"}},
		{"digest.keyword": Map{"order": "asc", "unmapped_type": "boolean"}},
	},
	"aggs": Map{
		"all": Map{
			"terms": Map{"field": "owner_id", "size": 10000},
			"aggs": Map{
				"edge": Map{
					"terms": Map{"field": "edge_serial.keyword", "size": 10000},
					"aggs": Map{
						"users": Map{
							"terms": Map{"field": "user_name.keyword", "size": 10000},
							"aggs": Map{
								"apps_cat": Map{
									"terms": Map{"field": "category.application", "size": 10000},
									"aggs": Map{
										"apps": Map{
											"terms": Map{"field": "detected_application_name.keyword", "size": 10000},
											"aggs": Map{
												"domains": Map{
													"terms": Map{"field": "host_server_name.keyword", "missing": "Unknown", "size": 10000},
													"aggs": Map{
														"duration": Map{
															"sum": Map{
																"script": Map{"source": "doc['last_seen_at'].value - doc['first_seen_at'].value"},
															},
														},
														"total_bytes": Map{
															"sum": Map{
																"field": "total_bytes",
															},
														},
													},
												},
												"traffic_bytes": Map{
													"sum": Map{
														"field": "total_bytes",
													},
												},
												"max_duration": Map{
													"max_bucket": Map{
														"buckets_path": "domains>duration",
													},
												},
											},
										},
									},
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
			"must_not": []Map{
				{"term": Map{"detected_protocol": 5}},
				{"term": Map{"detected_protocol": 81}},
			},
		},
	},
}
