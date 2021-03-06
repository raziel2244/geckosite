package model

import (
	"context"
	"database/sql"

	"github.com/eth0net/geckosite/s3"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// Animal stores the details for an animal.
type Animal struct {
	UUID
	Common

	Name        sql.NullString `json:"name" gorm:"unique"`
	Reference   sql.NullString `json:"reference" gorm:"unique"`
	Description sql.NullString `json:"description"`
	Sex         string         `json:"sex" gorm:"default:Unknown;not null;check:sex IN ('Male','Female','Unknown')"`
	Status      string         `json:"status" gorm:"default:Holdback;not null;check:status IN ('Non-Breeder','Breeder','Future Breeder','Holdback','For Sale','Sold','Other')"`
	Species     *Species       `json:"species"`
	SpeciesID   *uuid.UUID     `json:"-" gorm:"type:uuid;not null"`
	LaidAt      sql.NullTime   `json:"laidAt" gorm:"type:date"`
	HatchedAt   sql.NullTime   `json:"hatchedAt" gorm:"type:date"`

	BoughtAt     sql.NullTime `json:"boughtAt" gorm:"type:date"`
	BoughtFrom   *Contact     `json:"boughtFrom"`
	BoughtFromID *uuid.UUID   `json:"-" gorm:"type:uuid"`
	SoldAt       sql.NullTime `json:"soldAt" gorm:"type:date"`
	SoldTo       *Contact     `json:"soldTo"`
	SoldToID     *uuid.UUID   `json:"-" gorm:"type:uuid"`

	Parents  []*Animal `json:"parents" gorm:"many2many:animal_parents;foreignKey:id;joinForeignKey:child_id;references:id,sex;joinReferences:parent_id,parent_sex"`
	Children []*Animal `json:"children" gorm:"many2many:animal_parents;foreignKey:id,sex;joinForeignKey:parent_id,parent_sex;references:id;joinReferences:child_id"`

	Measurements []*Measurement `json:"measurements"`
	Traits       []*Trait       `json:"traits" gorm:"many2many:animal_traits"`
	Transactions []*Transaction `json:"transactions" gorm:"many2many:transaction_animals"`
}

// Father returns the father of the animal, if it is certain.
// See Animal.PossibleFathers for a list of possible fathers.
func (a Animal) Father() (father *Animal) {
	possibleFathers := a.PossibleFathers()
	if len(possibleFathers) != 1 {
		return nil
	}
	return possibleFathers[0]
}

// PossibleFathers returns a list of all possible fathers for the animal.
// See Animal.Father for a certain single result, if known.
func (a Animal) PossibleFathers() (possibleFathers []*Animal) {
	for _, parent := range a.Parents {
		if parent.Sex != "Male" {
			continue
		}
		possibleFathers = append(possibleFathers, parent)
	}
	return possibleFathers
}

// Mother returns the mother of the animal, if it is certain.
// See Animal.PossibleMothers for a list of possible mothers.
func (a Animal) Mother() (mother *Animal) {
	possibleMothers := a.PossibleMothers()
	if len(possibleMothers) != 1 {
		return nil
	}
	return possibleMothers[0]
}

// PossibleMothers returns a list of all possible mothers for the animal.
// See Animal.Mother for a certain single result, if known.
func (a Animal) PossibleMothers() (possibleMothers []*Animal) {
	for _, parent := range a.Parents {
		if parent.Sex != "Female" {
			continue
		}
		possibleMothers = append(possibleMothers, parent)
	}
	return possibleMothers
}

// Length returns the most recent length measurement of the animal.
func (a Animal) Length() (length *Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Length" {
			continue
		}
		length = measurement
	}
	return length
}

// Lengths returns all of the length measurements for the animal.
func (a Animal) Lengths() (lengths []*Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Length" {
			continue
		}
		lengths = append(lengths, measurement)
	}
	return lengths
}

// Weight returns the most recent weight measurement of the animal.
func (a Animal) Weight() (weight *Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Weight" {
			continue
		}
		weight = measurement
	}
	return weight
}

// Weights returns all of the weight measurements for the animal.
func (a Animal) Weights() (weights []*Measurement) {
	for _, measurement := range a.Measurements {
		if measurement.Type != "Weight" {
			continue
		}
		weights = append(weights, measurement)
	}
	return weights
}

// Images returns a list of images for the animal.
func (a Animal) Images() (images []string) {
	ch := s3.Client.ListObjects(
		context.Background(),
		a.Species.Order,
		minio.ListObjectsOptions{
			Prefix:    a.Species.Type + "/" + a.ID.String(),
			Recursive: true,
		},
	)
	for object := range ch {
		path := "/s3/" + a.Species.Order + "/" + object.Key
		images = append(images, path)
	}
	return images
}
