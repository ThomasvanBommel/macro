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
	UserID   int    `json:"user_id"`
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

type Meal struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Entry struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	FoodID    int     `json:"food_id"`
	MealID    Meal    `json:"meal"`
	Date      string  `json:"log_date"`
	Servings  float64 `json:"servings"`
	CreatedAt string  `json:"created_at"`
}

type EntryRequest struct {
	UserID   int     `json:"user_id"`
	Date	 string  `json:"date"`
}