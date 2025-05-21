package structsandjson

import (
	"bytes"
	"fmt"
)
import "encoding/json"

// Student Creating structure
type Student struct {
	name   string
	branch string
	year   int
	id     int
}

// Teacher Creating nested structure
type Teacher struct {
	Name    string
	nation  string
	id      int
	details Student
}

// StudentComplex Student struct with an anonymous structure and fields
type StudentComplex struct {
	name    string
	details struct { // Anonymous inner structure for personal details
		enrollment int
		GPA        float64 // Standard field
	}
}

// Person simple struct
type Person struct {
	Name string
	Age  int
}

// Weapon represents a generic weapon in the game or system.
type Weapon struct {
	Name          string   // Name of the weapon (e.g., "Longsword", "Shortbow")
	Type          string   // Type of weapon (e.g., "Sword", "Bow", "Axe", "Gun")
	Damage        float64  // Amount of damage the weapon inflicts
	Weight        float64  // Weight of the weapon
	Durability    int      // Current durability of the weapon (if applicable)
	MaxDurability int      // Maximum durability of the weapon (if applicable)
	Properties    []string // Special properties (e.g., "Fire", "Poison", "Two-Handed")
	Rarity        string   // Rarity of the weapon (e.g., "Common", "Rare", "Epic")
}

func PlaygroundStructsAndJSON() {

	student := Student{name: "Pep", branch: "master", year: 2025, id: 12}
	teacher := Teacher{Name: "Tiang", nation: "chinese", id: 1, details: student}

	fmt.Println("teacher Name:", teacher.nation)
	fmt.Println("student Name:", student.name)

	fmt.Printf("Teacher %v:\n", teacher)

	studentComplex := StudentComplex{
		name: "Alice",
		details: struct {
			enrollment int
			GPA        float64
		}{ // Initialize the anonymous struct here
			enrollment: 101,
			GPA:        3.85,
		},
	}

	// You can also omit the field names if you provide values in the order they are declared:
	studentComplex2 := StudentComplex{
		"Bob Johnson",
		struct {
			enrollment int
			GPA        float64
		}{ // Initialize the anonymous struct here
			205,
			3.5,
		},
	}

	fmt.Println("student complex name:", studentComplex.name)
	fmt.Printf("student complex %v:\n", studentComplex)
	fmt.Printf("student complex 2 %v:\n", studentComplex2)

	person := Person{Name: "John", Age: 30}

	jsonData, _ := json.Marshal(&person)
	fmt.Printf("person json: %s\n", string(jsonData))

	jsonData, _ = json.Marshal(&teacher)
	jsonDataLower := bytes.ToLower(jsonData)
	fmt.Printf("Teaches in JSON, but only show the fields exposed (first letter uppercase) jsonData: %s\n", string(jsonDataLower))

	// Example of how to create and initialize a Weapon variable
	sword := Weapon{
		Name:          "Steel Sword",
		Type:          "Sword",
		Damage:        15.5,
		Weight:        3.2,
		Durability:    100,
		MaxDurability: 100,
		Properties:    []string{"Sharp", "Double"},
		Rarity:        "Common",
	}

	// Accessing struct fields
	fmt.Println("Name:", sword.Name)
	fmt.Println("Type:", sword.Type)
	fmt.Println("Damage:", sword.Damage)
	fmt.Println("Properties:", sword.Properties)

	// Modifying a field
	sword.Damage = 18.0
	fmt.Println("New Damage:", sword.Damage)
}
