package models

import ()

type FoodRecordAggregate struct {
	Fat       float64  `json:"fat"`
	Carbs     float64  `json:"carbs"`
	Protein   float64  `json:"protein"`
	Calories  float64  `json:"calories"`
	Breakfast MealData `json:"breakfast"`
	Lunch     MealData `json:"lunch"`
	Dinner    MealData `json:"dinner"`
	Snack     MealData `json:"snack"`
}

type MealData struct {
	Count    int64   `json:"count"`
	Calories float64 `json:"calories"`
}
