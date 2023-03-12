package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for course -file
type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"` //defining the author using the pointer not using any instance.

}

type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

// helper , middleware -file
func IsEmpty(c *Course) bool {
	// return c.CourseId == ""&& c.CourseName==""
	return c.CourseName == ""
}

//helper, check if copy course exists:change of plans.!!
// func IsCopy(c *Course) bool {
// 	return
// }

// fake DB
var courses []Course

func main() {
	//create update delete courses.
	//having slice for fake database.
	//gorilla-mux for routing correctly.

	fmt.Println("APIs bro APIs")
	r := mux.NewRouter()

	//seeding
	//inject some data in courses which is slice of type course.
	//fillup all properties.
	courses = append(courses, Course{CourseId: "7", CourseName: "reactJs", CoursePrice: 300, Author: &Author{FullName: "GauravJha", Website: "lazyflash.net"}})

	courses = append(courses, Course{CourseId: "9", CourseName: "Mern", CoursePrice: 300, Author: &Author{FullName: "GauravJha", Website: "lazy-flash.net"}})

	//routing:
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")    //not expecting values GET is fine.
	r.HandleFunc("/course", createOneCourse).Methods("POST")     //will be sending data.(bring me data)
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT") //bring me data as well as id.
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")

	//ugly listen to a port
	log.Fatal(http.ListenAndServe(":4000", r))
}

//controllers -file

//serve home route

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>its working</h1>"))
}

// another route
func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all courses")
	//how to set (explict) header.
	w.Header().Set("Content-Type", "application/json")
	//setting a method and type of values accepting.

	//how to throw all data from db (fake)
	//json new writer.
	json.NewEncoder(w).Encode(courses)
	//courses will be treated as json value. throw back to whoever making request.

}

// creating new func
func getOneCourse(w http.ResponseWriter, r *http.Request) {
	//take unique id from request.
	//then compare from slice using loop
	fmt.Println("get one course")
	w.Header().Set("Content-Type", "application/json")

	//grab id form request
	params := mux.Vars(r) //param will have property of id

	//loop thorough courses, find matching id and return the response.
	for _, course := range courses {
		if course.CourseId == params["id"] { //not yet routing done.
			json.NewEncoder(w).Encode(course.Author)
			return
		}
	}
	json.NewEncoder(w).Encode("no course found with given id")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	//sofar encoding json
	// now decode
	fmt.Println("create one course")
	w.Header().Set("Content-Type", "application/json")

	//what if: body is empty
	if r.Body == nil { //check
		json.NewEncoder(w).Encode("Please send some data")
	}

	//what about -{}
	var course Course
	//decode the value
	//json handling: two method
	//decode acc to stuct
	//or looping through
	_ = json.NewDecoder(r.Body).Decode(&course) //taking values in course
	//bc course has been passed as reference.
	if IsEmpty(&course) {
		json.NewEncoder(w).Encode("No data inside JSON")
		return
	}

	//todo check only if title is duplicate
	//loop, title matches with course.coursename,json

	// if IsCopy(&course){
	// 	json.NewEncoder(w).Encode("course exists!")
	// 	return
	// }
	// using for loop not creating new function

	//creating param
	// params := mux.Vars(r)
	for _, first := range courses {
		if first.CourseName == course.CourseName {
			json.NewEncoder(w).Encode("this course exist already")
			return
		}
	}

	//generate unique id, string
	//append course into courses

	rand.Seed(time.Now().UnixNano())               //unique number
	course.CourseId = strconv.Itoa(rand.Intn(100)) //rand.Intn(100)returns int, but string needed.
	courses = append(courses, course)              //append values into courses
	json.NewEncoder(w).Encode(course)
	return

}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	//resend all value present or not.
	// loop through value find and update
	//remaking slice
	fmt.Println("update one course")
	w.Header().Set("Content-Type", "application/json")

	//first -grab id form req
	param := mux.Vars(r)
	//one value form url itself other one from the body.

	//loop(through the value), id , remmove, add(value again back in course) with my id
	for index, course := range courses {
		if course.CourseId == param["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			//add new values
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = param["id"]     //overriding course id if it has
			courses = append(courses, course) //appended course to list
			json.NewEncoder(w).Encode(course) //encoded to json.
			return
		}
	}
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete one course")
	w.Header().Set("Content-Type", "application/json")
	//find which course take id
	//loop through and find
	//use slices and index to remove.

	params := mux.Vars(r) //take id form request

	//loop,id,remove(index,index+1)

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			//TODO: --Send a confirm or deny response.
			json.NewEncoder(w).Encode("this id is now deleted")
			break
		}
	}
}

//todo Handling routes and testing routes
//going in main now.
