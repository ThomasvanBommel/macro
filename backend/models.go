package main

import "encoding/json"

type APIResponse interface {
	Status() int
	Result() any
}

type ErrorResponse struct {
	Code  int
	Quiet bool
	error
}

func (r *ErrorResponse) Status() int { return r.Code }
func (r *ErrorResponse) Result() any { return r.error }

type DataResponse struct {
	Code int
	Data any
}

func (r *DataResponse) Status() int { return r.Code }
func (r *DataResponse) Result() any { return r.Data }

type API struct {
	db *Database

	// 2xx Successful responses
	OK             func(any) APIResponse
	Message        func(string) APIResponse
	Created        func(any) APIResponse
	Accepted       func() APIResponse
	NoContent      func() APIResponse
	PartialContent func(any) APIResponse

	// 4xx Client error responses
	BadRequest      func(error) APIResponse
	Unauthorized    func(error) APIResponse
	Forbidden       func(error) APIResponse
	NotFound        func(error) APIResponse
	Conflict        func(error) APIResponse
	TooManyRequests func(error) APIResponse

	// 5xx Server error responses
	InternalServerError func(error) APIResponse
	NotImplemented      func(error) APIResponse
}

// API ---------------------------------------------------------------------------------------------

type UserCredentialInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ListUserEntriesInput struct {
	Name string `json:"name" binding:"required"`
	Date string `json:"date" binding:"required"`
}

type CreateFoodInput struct {
	Name         string      `json:"name" binding:"required"`
	Brand        string      `json:"brand"`
	Calories     json.Number `json:"calories" binding:"required"`
	Carbs        json.Number `json:"carbs" binding:"required"`
	Protein      json.Number `json:"protein" binding:"required"`
	Fat          json.Number `json:"fat" binding:"required"`
	ServingSize  string      `json:"serving_size" binding:"required"`
	ServingCount json.Number `json:"serving_count" binding:"required"`
}

type CreateEntryInput struct {
	FoodId   *int        `json:"food_id" binding:"required"`
	MealName string      `json:"meal_name" binding:"required"`
	Date     string      `json:"date" binding:"required"`
	Servings json.Number `json:"servings" binding:"required"`
}

type FoodResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Brand        string  `json:"brand"`
	Created      string  `json:"created"`
	UserName     string  `json:"username"`
	Calories     float64 `json:"calories"`
	Carbs        float64 `json:"carbs"`
	Protein      float64 `json:"protein"`
	Fat          float64 `json:"fat"`
	ServingSize  string  `json:"serving_size"`
	ServingCount float64 `json:"serving_count"`
}

type EntryWithFoodResponse struct {
	ID       int          `json:"id"`
	UserName string       `json:"username"`
	Food     FoodResponse `json:"food"`
	MealName string       `json:"meal_name"`
	Date     string       `json:"date"`
	Servings float64      `json:"servings"`
	Created  string       `json:"created"`
}

// Database ----------------------------------------------------------------------------------------

type User struct {
	Name         string
	PasswordHash string
	Created      string
}

type Session struct {
	UserName string `json:"username"`
	Token    string `json:"token"`
	Created  string `json:"created"`
	Expires  string `json:"expires"`
}

type Entry struct {
	ID       int
	UserName string
	FoodId   int
	MealName string
	Date     string
	Servings int
	Created  string
}

type Food struct {
	ID           int
	Name         string
	Brand        string
	Created      string
	UserName     string
	Calories     int
	Carbs        int
	Protein      int
	Fat          int
	ServingSize  string
	ServingCount int
}

type EntryWithFood struct {
	ID       int
	UserName string
	Food     Food
	MealName string
	Date     string
	Servings int
	Created  string
}

type CreateFoodParams struct {
	Name         string
	Brand        string
	Calories     int
	Carbs        int
	Protein      int
	Fat          int
	ServingSize  string
	ServingCount int
}

type EntryParams struct {
	FoodId   int
	MealName string
	Date     string
	Servings int
}
