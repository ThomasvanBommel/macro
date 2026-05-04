package db2

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
	ServingCount int
	ServingSize  string
}

type FoodParams struct {
	Name         string
	Brand        string
	Calories     int
	Carbs        int
	Protein      int
	Fat          int
	ServingCount int
	ServingSize  string
}
