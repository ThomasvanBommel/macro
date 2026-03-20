package main

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Session struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt string `json:"expires_at"`
}

type Food struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Brand       string `json:"brand,omitempty"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   int    `json:"created_by"`
	Calories    int    `json:"calories"`
	Carbs       int    `json:"carbs"`
	Fat         int    `json:"fat"`
	Protein     int    `json:"protein"`
	ServingSize int    `json:"serving_size"`
}

type PutFoodRequest struct {
	Name         string `json:"name" binding:"required"`
	Brand        string `json:"brand"`
	Calories     int    `json:"calories" binding:"required"`
	Carbs        int    `json:"carbs" binding:"required"`
	Fat          int    `json:"fat" binding:"required"`
	Protein      int    `json:"protein" binding:"required"`
	ServingSize  int    `json:"serving_size" binding:"required"`
	SessionToken string `json:"session_token"`
}

type Meal struct {
	Name string `json:"name"`
}

type Entry struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	FoodID    int     `json:"food_id"`
	Meal      Meal    `json:"meal"`
	Date      string  `json:"log_date"`
	Servings  float64 `json:"servings"`
	CreatedAt string  `json:"created_at"`
}

type PutEntryRequest struct {
	FoodID	     int     `json:"food_id" binding:"required"`
	Meal         string  `json:"meal" binding:"required"`
	Date         string  `json:"date" binding:"required"`
	Servings     float64 `json:"servings" binding:"required"`
	SessionToken string  `json:"session_token"`
}

type EntryRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
}
