package main

type Booking struct {
	Service string
	Doctor  string
	Date    string
	Time    string
}

type UserState struct {
	Step    string
	Booking Booking
}

var users = make(map[int64]*UserState)