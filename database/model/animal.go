package model

import (
	"time"

	"github.com/google/uuid"
)

// Animal stores the details for an animal.
type Animal struct {
	UUID
	Common

	// Details specific to the animal.
	Reference   string     `json:"reference"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Images      []string   `json:"images" gorm:"-"`
	Species     *Species   `json:"species"`
	SpeciesID   *uuid.UUID `json:"-" gorm:"not null"`
	Sex         string     `json:"sex" gorm:"default:Unknown;not null;check:sex IN ('Male','Female','Unknown')"`
	Status      string     `json:"status" gorm:"not null;check:status IN ('Non-Breeder','Breeder','Future Breeder','Holdback','For Sale','Sold')"`

	// Important dates for our records.
	DateLaid    *time.Time `json:"dateLaid"`
	DateHatched *time.Time `json:"dateHatched"`
	DateBought  *time.Time `json:"dateBought"`
	DateSold    *time.Time `json:"dateSold"`

	// Relations to other animals.
	Father   *Animal    `json:"father"`
	FatherID *uuid.UUID `json:"-"`
	Mother   *Animal    `json:"mother"`
	MotherID *uuid.UUID `json:"-"`

	// PossibleFathers []*Animal `json:"possibleFathers" gorm:"many2many:animal_parents;foreignKey:id,sex;"`
	// PossibleMothers []*Animal `json:"possibleMothers" gorm:"many2many:animal_parents;foreignKey:id,sex;"`
	// Children []Animal

	// Traits         []*Gene `json:"traits"`
	// PossibleTraits []*Gene `json:"possibleTraits"`

	Measurements []*Measurement `json:"measurements"`
}

// Weight returns the most recent weight measurement of the animal.
func (a Animal) Weight() (weight *Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Weight" {
			continue
		}
		weight = measurement
	}
	return
}

// Weights returns all of the weight measurements for the animal.
func (a Animal) Weights() (weights []*Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Weight" {
			continue
		}
		weights = append(weights, measurement)
	}
	return
}

// Length returns the most recent length measurement of the animal.
func (a Animal) Length() (length *Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Length" {
			continue
		}
		length = measurement
	}
	return
}

// Lengths returns all of the length measurements for the animal.
func (a Animal) Lengths() (lengths []*Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Length" {
			continue
		}
		lengths = append(lengths, measurement)
	}
	return
}
