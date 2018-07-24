package xlog

type TrendID struct {
	Hour   int `json:"hour" bson:"hour"`
	Minute int `json:"minute" bson:"minute"`
}

type Trend struct {
	ID    TrendID `json:"_id" bson:"_id"`
	Count int     `json:"count" bson:"count"`
}
