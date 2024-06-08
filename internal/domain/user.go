package domain

import "time"

type User struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Password          string `json:"password"`
	Mail              string `json:"mail"`
	Phone             string `json:"phone"`
	HealthDescription string `json:"health_description"`
	Trainings         []int  `json:"trainings"`
}

type Trainer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
	Trainings []int  `json:"trainings"`
}

type Training struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Level     string    `json:"level"`
	TrainerID int       `json:"trainer_id"`
	StartTime time.Time `json:"start_time" swaggertype:"string" example:"2024-06-08T15:04:05Z"`
	EndTime   time.Time `json:"end_time" swaggertype:"string" example:"2024-06-08T16:04:05Z"`
	Users     []int     `json:"users"`
}
