package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Up!\n")
}

// activities ..
func retrieveActivities(w http.ResponseWriter, r *http.Request) {

	activities := []Activity{}

	var activitiesWrapper Activities

	log.Println("Retrieve activtities - before")
	activities, err := getActivtities(ctx)
	log.Println("Retrieve activtities - after")

	if err != nil {
		http.Error(w, "whoops - activities not found", 503)
	}

	// for _, act := range storedActivities {
	// 	activities = append(activities, act)
	// }

	activitiesWrapper.Activities = activities

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(activitiesWrapper); err != nil {
		panic(err)
	}
}

// retrieveComment ..
func retrieveComment(w http.ResponseWriter, r *http.Request) {
	retrieveActivity(w, r, "comment")
}

// retrieveBlog ..
func retrieveBlog(w http.ResponseWriter, r *http.Request) {
	retrieveActivity(w, r, "blog")
}

// retrieve ..
func retrieveActivity(w http.ResponseWriter, r *http.Request, activityType string) {
	params := mux.Vars(r)
	uid := params["actUUID"]

	var activity *Activity
	activity, err := getActivity(ctx, uid)

	if err != nil {
		fmt.Println("whoops/404")
		http.NotFound(w, r)
	}

	if activity == nil {
		fmt.Println("whoops/404")
		http.NotFound(w, r)
	}

	if (*activity).GetType() != activityType {
		http.Error(w, "Incorrect activity type", 403)
		return
	}

	returnActivity(w, r, *activity)

}

// newBlog ..
func newBlog(w http.ResponseWriter, r *http.Request) {
	var blogMsg BlogMsg
	var authBlog AuthenticatedBlog
	var blog Blog

	log.Println("i")
	if r.Body == nil {
		log.Println("Error : Empty body")
		http.Error(w, "Please send a blog request body", 400)
		return
	}

	log.Println("ii")
	log.Printf("body: %s", r.Body)
	err := json.NewDecoder(r.Body).Decode(&blogMsg)
	if err != nil {
		log.Printf("Error : %s", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	log.Println("iii")
	log.Printf("authBlog: %s", authBlog)
	err = authenticateBlog(blogMsg.Blog)
	if err != nil {
		log.Printf("Error : %s", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	log.Println("iv")
	meta, err := freshMeta()
	if err != nil {
		log.Printf("Error : %s", err.Error())
		http.Error(w, err.Error(), 403)
		return
	}

	blog = Blog{
		Author:       blogMsg.Blog.Author,
		Title:        blogMsg.Blog.Title,
		Msg:          blogMsg.Blog.Msg,
		ActivityType: "blog",
		Meta:         meta,
	}

	log.Println(blogMsg.Blog.Author)
	log.Println(blog.Msg)
	log.Println(blogMsg.Blog.Msg)

	err = validateBlog(blog)
	if err != nil {
		log.Printf("Error : %s", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	log.Println("Validated blog")

	// storedBlog, err := storeActivity("blog", blog)
	log.Println("Stored blog")

	_, err = storeActivity("blog", blog)
	if err != nil {
		http.Error(w, err.Error(), 504)
		return

	}
	retrieveActivities(w, r)

}

// newBlog ..
func updateBlog(w http.ResponseWriter, r *http.Request) {
	var authBlog AuthenticatedBlog

	params := mux.Vars(r)
	uid := params["actUUID"]

	if r.Body == nil {
		http.Error(w, "Please send a blog request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&authBlog)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = authenticateBlog(authBlog)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	storedBlog, err := getBlog(ctx, uid)

	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	if storedBlog.ActivityType != "blog" {
		http.Error(w, "Incorrect activity type", 403)
		return
	}

	var blog Blog

	// blog = &storedBlog

	storedBlog.Title = authBlog.Title
	storedBlog.Author = authBlog.Author
	storedBlog.Msg = authBlog.Msg

	fmt.Println(blog.Msg)

	err = validateBlog(blog)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// storedBlog.Activity = blog
	updatedBlog, err := updateActivity(storedBlog)
	if err != nil {
		returnActivity(w, r, *updatedBlog)

	}
}

// new ..
func newComment(w http.ResponseWriter, r *http.Request) {
	var msgBody CommentMsg

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&msgBody)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(msgBody)

	meta, err := freshMeta()
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}
	log.Println("New comment - 1")

	comment := Comment{
		Meta:         meta,
		Author:       msgBody.Comment.Author,
		Msg:          msgBody.Comment.Msg,
		ActivityType: "comment",
	}

	err = validateComment(comment)
	log.Println("New comment - 2")

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println("New comment - 3")

	_, err = storeActivity("comment", comment)
	log.Println("New comment - 4")

	if err != nil {
		http.Error(w, err.Error(), 504)
		return

	}
	retrieveActivities(w, r)
}

// new ..
func updateComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment

	params := mux.Vars(r)
	uid := params["actUUID"]

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	newActivity, err := getActivity(ctx, uid)

	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	if (*newActivity).GetType() != "comment" {
		http.Error(w, "Incorrect activity type", 403)
		return
	}

	newComment, ok := (*newActivity).(Comment)
	if !ok {
		http.Error(w, "Incorrect activity type", 403)
		return
	}

	err = validateComment(newComment)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	// newActivity = newComment

	_, err = updateActivity(newComment)
	if err != nil {
		http.Error(w, err.Error(), 504)
		return

	}
	retrieveActivities(w, r)
}

func returnActivity(w http.ResponseWriter, r *http.Request, activity Activity) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(activity); err != nil {
		panic(err)
	}
}
