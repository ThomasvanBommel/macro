package db2

import "time"

type Food struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Brand        string    `json:"brand"`
	Created      time.Time `json:"created"`
	UserName     string    `json:"user_name"`
	Calories     Decimal   `json:"calories"`
	Carbs        Decimal   `json:"carbs"`
	Protein      Decimal   `json:"protein"`
	Fat          Decimal   `json:"fat"`
	ServingCount Decimal   `json:"serving_count"`
	ServingSize  string    `json:"serving_size"`
}

type NewFood struct {
	Name         string  `json:"name" binding:"required,ascii,min=3,max=30"`
	Brand        string  `json:"brand" binding:"ascii,max=30"`
	UserName     string  `json:"user_name" binding:"required,ascii,max=25"`
	Calories     Decimal `json:"calories" binding:"required,gte=0"`
	Carbs        Decimal `json:"carbs" binding:"required,gte=0"`
	Protein      Decimal `json:"protein" binding:"required,gte=0"`
	Fat          Decimal `json:"fat" binding:"required,gte=0"`
	ServingCount Decimal `json:"serving_count" binding:"required,gte=0.01"`
	ServingSize  string  `json:"serving_size" binding:"required,ascii,max=10"`
}
