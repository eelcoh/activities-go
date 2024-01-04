package main

import (
	api "github.com/eelcoh/go-api"
)

var routes = api.Routes{

	api.NewRoute(
		"Get Activities",
		"GET",
		"/activities",
		retrieveActivities,
	),

	api.NewRoute(
		"Get Comment",
		"GET",
		"/activities/comments/{actUUID}",
		retrieveComment,
	),

	api.NewRoute(
		"Get Blog",
		"GET",
		"/activities/blogs/{actUUID}",
		retrieveBlog,
	),

	api.NewRoute(
		"Update Comment",
		"PUT",
		"/activities/comments/{actUUID}",
		updateComment,
	),

	api.NewRoute(
		"Update Activity",
		"PUT",
		"/activities/blogs/{actUUID}",
		updateBlog,
	),

	api.NewRoute(
		"New Comment",
		"POST",
		"/activities/comments",
		newComment,
	),

	api.NewRoute(
		"New Blog",
		"POST",
		"/activities/blogs",
		newBlog,
		// NewRoute requires http.HandlerFunc,
		// newBlog is http.HandlerFunc,
		// authMiddleware expects http.Handler, outputs http.Handler
	),

	api.NewRoute(
		"Health",
		"GET",
		"/health",
		health,
	),
}
