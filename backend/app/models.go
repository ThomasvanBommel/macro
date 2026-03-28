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
	UserName string `json:"user_name"`
	Created  string `json:"created"`
	Expires  string `json:"expires"`
}

type ResponseFoodModel struct {
	ID           int    `json:"id"`
	Name         string `json:"name" binding:"required"`
	Brand        string `json:"brand"`
	Created      string `json:"created"`
	UserName     string `json:"user_name"`
	Calories     int    `json:"calories" binding:"gte=0"`
	Carbs        int    `json:"carbs" binding:"gte=0"`
	Protein      int    `json:"protein" binding:"gte=0"`
	Fat          int    `json:"fat" binding:"gte=0"`
	ServingSize  string `json:"serving_size" binding:"required"`
	ServingCount int    `json:"serving_count" binding:"gte=0"`
}

type RequestFoodModel struct {
	ResponseFoodModel
}

type RequestCreateFoodModel struct {
	Name         string `json:"name" binding:"required"`
	Brand        string `json:"brand"`
	Calories     int    `json:"calories" binding:"gte=0"`
	Carbs        int    `json:"carbs" binding:"gte=0"`
	Protein      int    `json:"protein" binding:"gte=0"`
	Fat          int    `json:"fat" binding:"gte=0"`
	ServingSize  string `json:"serving_size" binding:"required"`
	ServingCount int    `json:"serving_count" binding:"gte=0"`
}

type RequestEntriesModel struct {
	UserName string `json:"user_name" binding:"required"`
	Date     string `json:"date" binding:"required"`
}

// type RequestAddEntryModel struct {
// 	FoodId   *int              `json:"food_id"`
// 	Food     *RequestFoodModel `json:"food"`
// 	MealName string            `json:"meal_name" binding:"required"`
// 	Date     string            `json:"date" binding:"required"`
// 	Servings int               `json:"servings" binding:"gte=1"`
// }

type ResponseEntryModel struct {
	ID       int               `json:"id"`
	UserName string            `json:"user_name"`
	Food     ResponseFoodModel `json:"food"`
	MealName string            `json:"meal_name"`
	Date     string            `json:"date"`
	Servings int               `json:"servings"`
	Created  string            `json:"created"`
}

type RequestAddEntryModel struct {
	FoodId   *int   `json:"food_id" binding:"required"`
	MealName string `json:"meal_name" binding:"required"`
	Date     string `json:"date" binding:"required"`
	Servings int    `json:"servings" binding:"gt=0"`
}
