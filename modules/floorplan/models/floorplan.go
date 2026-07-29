package models

type FloorTable struct {
	Id       string  `json:"id" bson:"id"`
	Label    string  `json:"label" bson:"label"`
	Zone     string  `json:"zone" bson:"zone"`
	Capacity int     `json:"capacity" bson:"capacity"`
	Shape    string  `json:"shape" bson:"shape"`
	X        float64 `json:"x" bson:"x"`
	Y        float64 `json:"y" bson:"y"`
	Width    float64 `json:"width" bson:"width"`
	Height   float64 `json:"height" bson:"height"`
	Status   string  `json:"status" bson:"status"`
}

type FloorZone struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	W    int    `json:"w" bson:"w"`
	H    int    `json:"h" bson:"h"`
}
