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
	test5 := "../Data/test5.txt"
	//test10 := "../Data/test10.txt"
	//test20 := "../Data/test20.txt"
	updatedTestData := userBasedPearsonPrediction(test5)
	testData := getTestData(test5)
	fmt.Println(getRMSE(testData, updatedTestData)) //prints the RMSE to evaluate the prediction
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
	data, err := ioutil.ReadFile(filename) //read contents of text file into data array
	if err != nil {
		fmt.Println("Failed to read test file")
		return [100][1000]int{}
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data))) //scan input as string of space-delimited tokens
	scanner.Split(bufio.ScanLines)                               //split function for only scanning words, not spaces
	checkUsers := []int{}                                        //slice to store userid's for checking
	var checker int
	users := [100][1000]int{} //100 users starting at userid 0, 1000 movies starting at movieid 0
	//ensure only the movies for each user as given by the test data applies
	for i := 0; i < 100; i++ {
		for j := 0; j < 1000; j++ {
			users[i][j] = 9
		}
	}
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
			for i := 0; i < len(checkUsers); i++ { //traverse slice to ensure there are no duplicates for userid
				if checkUsers[i] == userId {
					checker = 1
				}
			}
			if checker == 0 { //adding new userid to users array
				checkUsers = append(checkUsers, userId)
			}
		}
		//populate 2d array
		users[len(checkUsers)-1][movieId-1] = rating
	}
	return users
}

//Function to find the pearson similarity between two users
func userPearson(userOne int, userTwo int, filename string) float64 {
	trainingData := getData() //Get the training data
	testData := getTestData(filename)
	train := "train"
	var numerator int = 0
	var denominator float64 = 0
	var squaredOne int = 0
	var squaredTwo int = 0
	var similarity float64 = 0
	activeUserAverage := getAverageRating(userOne, filename)
	otherUserAverage := getAverageRating(userTwo, train)

	//traverse through the movie set to find the movies userOne and userTwo have rated
	for movieId := 0; movieId < 1000; movieId++ {
		if testData[userOne][movieId] != 0 && testData[userOne][movieId] != 9 && trainingData[userTwo][movieId] != 0 { //both users have rated the same movie
			//determine similarity score between users
			numerator = numerator + ((testData[userOne][movieId] - activeUserAverage) * (trainingData[userTwo][movieId] - otherUserAverage))
			squaredOne = squaredOne + ((testData[userOne][movieId] - activeUserAverage) * (testData[userOne][movieId] - activeUserAverage))
			squaredTwo = squaredTwo + ((trainingData[userTwo][movieId] - otherUserAverage) * (trainingData[userTwo][movieId] - otherUserAverage))
		}
	}
	denominator = math.Sqrt(float64(squaredOne)) * math.Sqrt(float64(squaredTwo))
	similarity = float64(numerator) / denominator
	return similarity
}

func getAverageRating(user int, filename string) int {
	var sum int
	var num int

	//for the other user
	if filename == "train" {
		trainData := getData()
		for movie := 0; movie < 1000; movie++ {
			if trainData[user][movie] != 0 { //user has rated that movie
				sum += trainData[user][movie]
				num++
			}
		}
		return int(float64(sum) / float64(num)) //returns average rating for given user
	} else { //for the active user
		testData := getTestData(filename)
		for movie := 0; movie < 1000; movie++ {
			if testData[user][movie] != 0 && testData[user][movie] != 9 { //user has rated that movie
				sum += testData[user][movie]
				num++
			}
		}
		return int(float64(sum) / float64(num))
	}
}

func getPrediction(activeUser int, activeMovie int, filename string) int {
	trainingData := getData() //get training data for similarity prediction
	train := "train"
	kSimilarUsers := [15]int{}       //array to store the top 15 most similar users
	kSimilarRatings := [15]float64{} //array to store similarity scores for the top 15 most similar users
	leastSimilar := 0                //int variable to keep track of the index for the similar user in kSimilarUsers with the smallest similarity score

	//find k most similar users to active user given active movie
	for other := 0; other < 200; other++ {
		if trainingData[other][activeMovie] != 0 { //other user has also rated the same movie
			similarityScore := userPearson(activeUser, other, filename) //returns similarity score between active user and other user
			if math.Abs(similarityScore) > math.Abs(kSimilarRatings[leastSimilar]) {
				kSimilarUsers[leastSimilar] = other             //update most similar users array
				kSimilarRatings[leastSimilar] = similarityScore //update most similar users' ratings array
			}
			//traverse the similar users and ratings arrays to find the least similar user
			for i := 0; i < 15; i++ {
				if math.Abs(kSimilarRatings[i]) < math.Abs(kSimilarRatings[leastSimilar]) {
					leastSimilar = i
				}
			}
		}
	}

	//predict ratings for active user given most similar 15 users
	var numerator float64
	var denominator float64
	var prediction float64
	activeUserAverage := getAverageRating(activeUser, filename)

	for u := 0; u < 15; u++ { //traverse through the top 15 most similar users
		otherUserAverage := getAverageRating(u, train)
		numerator = numerator + (kSimilarRatings[u] * (float64(trainingData[kSimilarUsers[u]][activeMovie] - otherUserAverage)))
		denominator = denominator + math.Abs(kSimilarRatings[u])
	}

	prediction = float64(activeUserAverage) + (math.Abs(numerator) / denominator)
	return int(math.Abs(math.Round(prediction)))
}

//Function implements user-based collaborative filtering algorithm with cosine similarity
func userBasedPearsonPrediction(filename string) [100][1000]int {
	testData := getTestData(filename)
	updatedTestData := testData
	counter := 0

	for user := 0; user < 100; user++ {
		for movie := 0; movie < 1000; movie++ {
			if testData[user][movie] == 0 { //need prediction for active user
				prediction := getPrediction(user, movie, filename) //to get prediction for movie item for active user
				updatedTestData[user][movie] = prediction
			}
		}
		counter++
		fmt.Println("Prediction done for user:", counter)
	}
	return updatedTestData
}

//Function to evaluate prediction
func getRMSE(testData [100][1000]int, predictedData [100][1000]int) float64 {
	fmt.Print("Evaluating predictions")
	actualData := getData() //retrieves the trainingData
	var totalSum int = 0
	var totalMissing int = 0
	//traverse the missing ratings in the test set
	for user := 0; user < 100; user++ {
		for movie := 0; movie < 100; movie++ {
			if testData[user][movie] == 0 { //denotes the missing ratings
				totalMissing++
				totalSum = totalSum + ((predictedData[user][movie] - actualData[user][movie]) * (predictedData[user][movie] - actualData[user][movie]))
			}
		}
	}
	score := math.Sqrt(float64(totalSum / totalMissing))
	return score
}
