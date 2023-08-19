package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Person struct {
	Name   string
	Score  int
	Pref1  int
	Pref2  int
}

type Player {
	UserName string
	DisplayName string
	Type string
	Build string
	EHP	float64
	EHB float64
}

type Membership{
	PlayerId float64
	Player Player
}

type Group struct {
	id int
	Name string
	ClanChat string
	Memberships []Membership
}

// Sample data (replace this with actual data)
// Each person has a weighted score and two preference values (pref_1 and pref_2)
var people = []Person{
	{"Alice", 90, 5, 3},
	{"Bob", 85, 2, 5},
	{"Charlie", 78, 3, 4},
	// Add more people here
}

// Constants and parameters
const (
	NUM_TEAMS      = 2
	POPULATION_SIZE = 100
	MUTATION_RATE   = 0.1
	GENERATIONS     = 100
)

// Function to calculate a person's overall score based on weighted score and preferences
func calculateOverallScore(person Person, weightScore, weightPref1, weightPref2 float64) float64 {
	return (float64(person.Score) * weightScore) + (float64(person.Pref1) * weightPref1) + (float64(person.Pref2) * weightPref2)
}

// Function to create a random team assignment
func createRandomTeamAssignment(people []Person) [][]Person {
	rand.Shuffle(len(people), func(i, j int) {
		people[i], people[j] = people[j], people[i]
	})

	teamAssignment := make([][]Person, NUM_TEAMS)
	for i := 0; i < len(people); i++ {
		teamAssignment[i%NUM_TEAMS] = append(teamAssignment[i%NUM_TEAMS], people[i])
	}

	return teamAssignment
}

// Function to evaluate the fitness of a team assignment (lower value is better)
func evaluateFitness(teamAssignment [][]Person, weightScore, weightPref1, weightPref2 float64) float64 {
	teamScores := make([]float64, NUM_TEAMS)
	for i, team := range teamAssignment {
		for _, person := range team {
			teamScores[i] += calculateOverallScore(person, weightScore, weightPref1, weightPref2)
		}
	}
	sort.Float64s(teamScores)
	return teamScores[NUM_TEAMS-1] - teamScores[0]
}

// Function to perform a single mutation on the team assignment
func mutateTeamAssignment(teamAssignment [][]Person) {
	i := rand.Intn(NUM_TEAMS)
	j := (i + rand.Intn(NUM_TEAMS-1) + 1) % NUM_TEAMS
	teamAssignment[i], teamAssignment[j] = teamAssignment[j], teamAssignment[i]
}

// Genetic Algorithm to fairly distribute teams
func geneticAlgorithm(people []Person, weightScore, weightPref1, weightPref2 float64) [][]Person {
	population := make([][][]Person, POPULATION_SIZE)
	for i := 0; i < POPULATION_SIZE; i++ {
		population[i] = createRandomTeamAssignment(people)
	}

	for generation := 0; generation < GENERATIONS; generation++ {
		sort.SliceStable(population, func(i, j int) bool {
			return evaluateFitness(population[i], weightScore, weightPref1, weightPref2) < evaluateFitness(population[j], weightScore, weightPref1, weightPref2)
		})

		if generation%10 == 0 {
			fmt.Printf("Generation %d, Fitness: %.2f\n", generation, evaluateFitness(population[0], weightScore, weightPref1, weightPref2))
		}

		newPopulation := population[:POPULATION_SIZE/2]

		for len(newPopulation) < POPULATION_SIZE {
			parent1 := population[rand.Intn(POPULATION_SIZE/2)]
			parent2 := population[rand.Intn(POPULATION_SIZE/2)]
			child := make([][]Person, NUM_TEAMS)
			for i := 0; i < NUM_TEAMS; i++ {
				child[i] = append(child[i], parent1[i]...)
			}
			mutateTeamAssignment(child)
			newPopulation = append(newPopulation, child)
		}

		population = newPopulation
	}

	return population[0]
}

func main() {
	// Weight values (you can adjust these based on your requirements)
	weightScore := 1.0
	weightPref1 := 0.5
	weightPref2 := 0.5

	// Run the genetic algorithm
	finalTeamAssignment := geneticAlgorithm(people, weightScore, weightPref1, weightPref2)

	// Display the final team assignment
	for i, team := range finalTeamAssignment {
		sort.SliceStable(team, func(j,k int) bool {
			return team[i][j].Score < team[i][k].Score
		})
		fmt.Printf("Team %d: ", i+1)
		for _, person := range team {
			fmt.Printf("%s ", person.Name)
		}
		fmt.Println()
	}
}
