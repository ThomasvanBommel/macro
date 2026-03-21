package main

type RequestUserModel struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResponseUserModel struct {
	Name    string `json:"name"`
	Created string `json:"created"`
}

type ResponseSessionModel struct {
	UserName  string `json:"user_name"`
	Created   string `json:"created"`
	Expires   string `json:"expires"`
}

/*
FoodModel:

	calories, carbs, protein, fat, and serving_count are stored as integer representing the value
	multiplied by 100 to avoid floating point issues.
*/
type FoodModel struct {
	id            int    `json:"id"`
	name          string `json:"name"`
	brand         string `json:"brand"`
	created       string `json:"created"`
	user_name     string `json:"username"`
	calories      int    `json:"calories"`
	carbs         int    `json:"carbs"`
	protein       int    `json:"protein"`
	fat           int    `json:"fat"`
	serving_size  string `json:"serving_size"`
	serving_count int    `json:"serving_count"`
}

type MealModel struct {
	name string `json:"name"`
}

/*
EntryModel:

	servings is stored as integer representing the value multiplied by 100 to avoid floating point
	issues.
*/
type EntryModel struct {
	id        int    `json:"id"`
	user_name string `json:"user_name"`
	food_id   int    `json:"food_id"`
	meal_name string `json:"meal_name"`
	date      string `json:"date"`
	servings  int    `json:"servings"`
	created   string `json:"created"`
}
