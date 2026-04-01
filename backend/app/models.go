package main

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
	Name         string `json:"name" binding:"required"`
	Brand        string `json:"brand"`
	Calories     int    `json:"calories" binding:"gte=0"`
	Carbs        int    `json:"carbs" binding:"gte=0"`
	Protein      int    `json:"protein" binding:"gte=0"`
	Fat          int    `json:"fat" binding:"gte=0"`
	ServingSize  string `json:"serving_size" binding:"required"`
	ServingCount int    `json:"serving_count" binding:"gt=0"`
}

type CreateEntryInput struct {
	FoodId   *int   `json:"food_id" binding:"required"`
	MealName string `json:"meal_name" binding:"required"`
	Date     string `json:"date" binding:"required"`
	Servings int    `json:"servings" binding:"gt=0"`
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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Brand        string `json:"brand"`
	Created      string `json:"created"`
	UserName     string `json:"username"`
	Calories     int    `json:"calories"`
	Carbs        int    `json:"carbs"`
	Protein      int    `json:"protein"`
	Fat          int    `json:"fat"`
	ServingSize  string `json:"serving_size"`
	ServingCount int    `json:"serving_count"`
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

type CreateEntryParams struct {
	FoodId   int
	MealName string
	Date     string
	Servings int
}
