package db

type User struct {
	Name         string
	PasswordHash string
	Created      string
}

type Session struct {
	UserName string
	Token    string
	Created  string
	Expires  string
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

type FoodParams struct {
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
