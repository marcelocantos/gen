package models

// +gen projections:"int,Thing2,string"
type Thing1 struct {
	Name   string
	Number int
}

// +gen methods:"Max,Min,Sort"
type Thing2 Thing3

type Thing3 float64
