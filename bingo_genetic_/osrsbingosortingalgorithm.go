package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"
)

type Player struct {
	UserName    string
	DisplayName string
	Type        string
	Build       string
	EHP         float64
	EHB         float64
}

type Membership struct {
	PlayerId float64
	Player   Player
}

type Group struct {
	id          int
	Name        string
	ClanChat    string
	Memberships []Membership
}

type Person struct {
	Name  string
	Score float64
	Pref  int
	Type  int
}

type PlayerPreferences struct {
	Name      string
	Pref      string
	BankValue int
}

var (
	playerTypeMap = map[string]int{
		"unknown":  0,
		"regular":  1,
		"ironman":  2,
		"hardcore": 3,
		"ultimate": 4,
	}
)

var (
	personTypeToplayerTypeMap = map[int]string{
		0: "gim",
		1: "regular",
		2: "ironman",
		3: "hardcore",
		4: "ultimate",
	}
)

const (
	NUM_TEAMS       = 10
	POPULATION_SIZE = 100
	MUTATION_RATE   = 0.1
	GENERATIONS     = 100
	EHB_MAX         = 200
	EHP_MAX         = 200
)

var CLOCK_ = time.Now().UnixNano()

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != (nil) {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
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

// TODO (8/4/2023): this function is currently weighted based on ehb and ehp which might need more fine tuning.
// Also add bank_value when that data is gathered
func calculateScore(player Player) float64 {
	var playerScore float64

	if player.EHB >= EHB_MAX {
		playerScore += 10.0
	} else {
		playerScore += (player.EHB * 10.0) / EHB_MAX
	}

	if player.EHP >= EHP_MAX {
		playerScore += 10.0
	} else {
		playerScore += (player.EHP * 10.0) / EHP_MAX
	}

	return playerScore
}

func getBingoPreference(Name string) int {
	//TODO(8/4/2023) : Collect player preference data
	// playerPrefMap := make(map[string]string)
	// switch playerPrefMap[Name] {
	// case "Skilling":
	// 	return 1
	// case "PVM":
	// 	return 2
	// }
	return (rand.Int() % 2) + 1
}

// Function to evaluate the fitness of a team assignment (lower value is better)
func evaluateFitness(teamAssignment [][]Person) float64 {
	teamScores := make([]float64, NUM_TEAMS)
	teamPreferences := make([]float64, NUM_TEAMS)

	for i, team := range teamAssignment {
		for _, person := range team {
			teamScores[i] += person.Score
			teamPreferences[i] += float64(person.Pref)
		}
	}

	// Sort by team score (lower value is better)
	sort.SliceStable(teamScores, func(i, j int) bool {
		return teamScores[i] < teamScores[j]
	})

	// Sort by team preference (lower value is better)
	sort.SliceStable(teamPreferences, func(i, j int) bool {
		return teamPreferences[i] < teamPreferences[j]
	})

	// Combine the scores and preferences for a final fitness evaluation
	// You can adjust the weight to prioritize score or preference more
	finalFitness := (teamScores[NUM_TEAMS-1] - teamScores[0]) + (teamPreferences[NUM_TEAMS-1] - teamPreferences[0])

	return finalFitness
}

// Function to perform a single mutation on the team assignment
func mutateTeamAssignment(teamAssignment [][]Person) {
	if rand.Float64() < MUTATION_RATE {
		i := rand.Intn(NUM_TEAMS)
		j := (i + rand.Intn(NUM_TEAMS)) % NUM_TEAMS
		individual1 := rand.Intn(len(teamAssignment[i]))
		individual2 := rand.Intn(len(teamAssignment[j]))
		teamAssignment[i][individual1], teamAssignment[j][individual2] = teamAssignment[j][individual2], teamAssignment[i][individual1]
	}
}

func loadBingoPreferenceData() [][]string {
	f, err := os.Open("responses.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	responses, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	return responses
}

func retrievePlayerData() Person {
	player := new(Player)
	getJson("https://api.wiseoldman.net/v2/players/", player)

	person := Person{
		Name:  player.DisplayName,
		Score: calculateScore(*player),
		Pref:  getBingoPreference(player.DisplayName),
		Type:  playerTypeMap[player.Type],
	}

	return person
}

func retrieveAndTransformGroupData() []Person {
	group := new(Group)
	getJson("https://api.wiseoldman.net/v2/groups/4329", group)

	people := []Person{}

	for i := 0; i < len(group.Memberships); i++ {
		player := group.Memberships[i].Player
		person := Person{
			Name:  player.DisplayName,
			Score: calculateScore(player),
			Pref:  getBingoPreference(player.DisplayName),
			Type:  playerTypeMap[player.Type],
		}
		people = append(people, person)
	}

	return people
}

func geneticAlgorithm(people []Person) [][]Person {
	population := make([][][]Person, POPULATION_SIZE)
	for i := 0; i < POPULATION_SIZE; i++ {
		population[i] = createRandomTeamAssignment(people)
	}

	for generation := 0; generation < GENERATIONS; generation++ {
		sort.SliceStable(population, func(i, j int) bool {
			return evaluateFitness(population[i]) < evaluateFitness(population[j])
		})
		if generation%10 == 0 {
			fmt.Printf("Generation %d, Fitness: %.2f\n", generation, evaluateFitness(population[0]))
		}

		newPopulation := population[:POPULATION_SIZE/2]
		for len(newPopulation) < POPULATION_SIZE {
			parent1 := population[rand.Intn(POPULATION_SIZE/2)]
			parent2 := population[rand.Intn(POPULATION_SIZE/2)]

			child := crossGenetics(parent1, parent2)

			mutateTeamAssignment(child)

			newPopulation = append(newPopulation, child)
		}
		population = newPopulation
	}

	sort.SliceStable(population, func(i, j int) bool {
		return evaluateFitness(population[i]) < evaluateFitness(population[j])
	})

	return population[0]
}

func crossGenetics(parent1, parent2 [][]Person) [][]Person {
	child := make([][]Person, NUM_TEAMS)

	for i := 0; i < NUM_TEAMS; i++ {
		// Create a new slice for the child team
		child[i] = make([]Person, len(parent1[i]))

		// Copy the people from parent1 to the child team
		copy(child[i], parent1[i])

		// Keep track of the individuals already added to the child team
		usedIndices := make(map[int]bool)

		// Randomly choose some unique people from parent2 to replace in the child team
		for j := 0; j < len(parent2[i]); j++ {
			// Find an available index in the child team
			index := rand.Intn(len(child[i]))
			for usedIndices[index] {
				index = rand.Intn(len(child[i]))
			}

			// Add the individual from parent2 to the child team
			child[i][index] = parent2[i][j]
			usedIndices[index] = true
		}
	}

	return child
}
func main() {
	// Seed the random number generator based on the current time
	rand.Seed(CLOCK_)
	people := retrieveAndTransformGroupData()
	teams := geneticAlgorithm(people)

	// Display the final team assignment
	for i, team := range teams {
		sort.SliceStable(team, func(j, k int) bool {
			return team[j].Score > team[k].Score
		})
		fmt.Printf("Team %d: \n", i+1)
		for _, person := range team {
			fmt.Printf("%s score: %.2f content_preference: %.d type: %s\n", person.Name, person.Score, person.Pref, personTypeToplayerTypeMap[person.Type])
		}
		fmt.Println()
	}
}
