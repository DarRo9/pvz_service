package repository

import "time"

type PVZ struct {
	ID               string    `db:"id"`
	City             string    `db:"city"`
	RegistrationDate time.Time `db:"registration_date"`
}

type Reception struct {
	ID            string    `db:"id"`
	ExecutionDate time.Time `db:"execution_date"`
	PVZID         string    `db:"pvz_id"`
	Status        string    `db:"status"`
}

type User struct {
	ID               string    `db:"id"`
	Email            string    `db:"email"`
	Password         string    `db:"password"`
	Role             string    `db:"role"`
	RegistrationDate time.Time `db:"registration_date"`
}

type Product struct {
	ID            string    `db:"id"`
	ReceptionDate time.Time `db:"reception_date"`
	ReceptionId   string    `db:"reception_id"`
	Type          string    `db:"type"`
}

type ReceptionWithProducts struct {
	Reception *Reception
	Products  []*Product
}

type PVZWithReceptions struct {
	PVZ        *PVZ
	Receptions []*ReceptionWithProducts
}
