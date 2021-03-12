//Declaring main package, groups functions and all files in same directory movieRecommendation
package main

//importing packages
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

//run file in terminal with "go run <filename>" command
func main() {
	fmt.Println(userBasedCosinePrediction("../Data/test5.txt"))
}

//Function to read the training data and insert data into 2d array (200 users x 1000 movies)
func getData() [200][1000]int { //func <function name> <returning value of specified type: 2d array of integers>
	data, err := ioutil.ReadFile("../Data/train.txt") //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read training file")
		return [200][1000]int{}
	}
	//scan input as sequence of space-delimited tokens
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Split(bufio.ScanWords) //split function for only scanning words, not spaces
	buffer := [200][1000]int{}     //buffer to store data in

	for row := 0; row < 200; row++ {
		for col := 0; col < 1000; col++ {
			scanner.Scan()                                   //advances scanner to next token
			currentRank, err := strconv.Atoi(scanner.Text()) //Atoi: string conversion to int
			if err != nil {
				fmt.Println("Error with scanner at token")
				return [200][1000]int{}
			}
			buffer[row][col] = currentRank //adding integer value pos into buffer at specified position
		}
	}
	return buffer
}

//Function to read the test data passed as an argument and insert data into 2d array (100 users x 1000 movies)
func getTestData(filename string) [100][1000]int {
	data, err := ioutil.ReadFile(filename) //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read test file")
		return [100][1000]int{}
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data))) //scan input as string of space-delimited tokens
	scanner.Split(bufio.ScanLines)                               //split function for only scanning words, not spaces
	checkUsers := []int{}
	var checker int
	users := [100][1000]int{} //100 users starting at userid 1, 1000 movies starting at movieid 1
	for scanner.Scan() {
		line := strings.Fields(scanner.Text()) //Fields function breaks a string around each instance of white space into an array
		currentUser := line[0]
		userId, _ := strconv.Atoi(currentUser) //Atoi: string conversion to int
		currentMovie := line[1]
		movieId, _ := strconv.Atoi(currentMovie)
		currentRating := line[2]
		rating, _ := strconv.Atoi(currentRating)
		//for first instance of test user
		if len(checkUsers) == 0 {
			checkUsers = append(checkUsers, userId)
		} else {
			checker = 0
			for i := 0; i < len(checkUsers); i++ {
				if checkUsers[i] == userId {
					checker = 1
				}
			}
			if checker == 0 {
				checkUsers = append(checkUsers, userId)
			}
		}
		//populate 2d array
		users[len(checkUsers)-1][movieId-1] = rating
	}
	return users
}

//Function to find the cosine similarity between two users
func userCosine(userOne int, userTwo int, filename string) float64 {
	trainingData := getData() //Get the training data
	testData := getTestData(filename)
	var numerator float64 = 0
	var denominator float64 = 0
	var squaredOne float64 = 0
	var squaredTwo float64 = 0
	var similarity float64 = 0

	//traverse through the movie set to find the movies userOne and userTwo have rated
	for movieId := 0; movieId < 1000; movieId++ {
		if testData[userOne][movieId] != 0 && trainingData[userTwo][movieId] != 0 { //both users have rated the same movie
			//determine similarity score between users
			numerator = numerator + (float64(trainingData[userOne][movieId]) * float64(trainingData[userTwo][movieId]))
			squaredOne = squaredOne + (float64(trainingData[userOne][movieId]) * float64(trainingData[userOne][movieId]))
			squaredTwo = squaredTwo + (float64(trainingData[userTwo][movieId]) * float64(trainingData[userTwo][movieId]))
		}
	}
	denominator = math.Sqrt(squaredOne) * math.Sqrt(squaredTwo)
	similarity = numerator / denominator
	return similarity
}

func getPrediction(activeUser int, activeMovie int, filename string) int {
	trainingData := getData() //get training data for similarity prediction
	kSimilarUsers := [15]int{} //array to store the top 15 most similar users
	kSimilarRatings := [15]float64{} //array to store similarity scores for the top 15 most similar users
	leastSimilar := 0 //int variable to keep track of the index for the similar user in kSimilarUsers with the smallest similarity score
	
	//find k most similar users to active user given active movie
	for other := 0; other < 200; other++ {
		if trainingData[other][activeMovie] != 0 { //other user has also rated the same movie
			similarityScore := userCosine(activeUser, other, filename) //returns similarity score between active user and other user

			if similarityScore > kSimilarRatings[leastSimilar] {
				kSimilarUsers[leastSimilar] = other //update most similar users array
				kSimilarRatings[leastSimilar] = similarityScore //update most similar users' ratings array
			}
			//traverse the similar users and ratings arrays to find the least similar user
			for i := 0; i < 15; i++ {
				if kSimilarRatings[i] < kSimilarRatings[leastSimilar] {
					leastSimilar = i
				}
			}
		}
	}
	
	//predict ratings for active user given most similar 15 users
	var numerator float64
	var denominator float64
	var prediction float64
	
	for u := 0; u < 15; u++ { //traverse through the top 15 most similar users
		numerator = numerator + (kSimilarRatings[u] * float64(trainingData[kSimilarUsers[u]][activeMovie]))
		denominator = denominator + kSimilarRatings[u]
	}

	prediction = numerator / denominator
	return int(math.Round(prediction))
} 

//Function implements user-based collaborative filtering algorithm with cosine similarity
func userBasedCosinePrediction(filename string)  [100][1000]int {
	testData := getTestData(filename)
	updatedTestData := testData

	for user := 0; user < 100; user++ {
		for movie := 0; movie < 1000; movie++ {
			if testData[user][movie] == 0 { //need prediction for active user
				prediction := getPrediction(user, movie, filename) //to get prediction for movie item for active user
				updatedTestData[user][movie] = prediction
			}
		}
	}
	return updatedTestData
}
