package main

type Booking struct {
	ServiceID int64
	DoctorID  int64
	Date      string
	Time      string
}

type UserState struct {
	Step        string
	Booking     Booking
	PatientName string
	Phone       string
}

var users = make(map[int64]*UserState)
