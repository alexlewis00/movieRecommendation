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
	fmt.Println(itemBasedCosinePrediction(test5))
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
	scanner.Split(bufio.ScanLines) //split function for only scanning words, not spaces
	checkUsers := []int{} //slice to store userid's for checking
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

//Function to find the adjusted cosine similarity between two users
func itemCosine(movieOne int, movieTwo int, filename string) float64 {
	trainingData := getData() //Get the training data
	testData := getTestData(filename)
	var numerator float64 = 0
	var denominator float64 = 0
	var squaredOne float64 = 0
	var squaredTwo float64 = 0
	var similarity float64 = 0

	totalUsers := []int{}
	for userId := 0; userId < 100; userId++ {
		if testData[userId][movieOne] != 0 && testData[userId][movieOne] != 9 && trainingData[userId][movieTwo] != 0 {
			totalUsers = append(totalUsers, userId)
		}
	}
	averageRating := getAverageRating(totalUsers, movieOne, filename)

	//traverse through the active user set to find the user that rated both movieOne and movieTwo have rated
	for userId := 0; userId < 100; userId++ {
		if testData[userId][movieOne] != 0 && testData[userId][movieOne] != 9 && trainingData[userId][movieTwo] != 0 { //both movies are rated
			//determine similarity score between movies
			numerator = numerator + (float64(trainingData[userId][movieOne] - averageRating) * float64(trainingData[userId][movieTwo] - averageRating))
			squaredOne = squaredOne + (float64(trainingData[userId][movieOne] - averageRating) * float64(trainingData[userId][movieOne] - averageRating))
			squaredTwo = squaredTwo + (float64(trainingData[userId][movieTwo] - averageRating) * float64(trainingData[userId][movieTwo] - averageRating))
		}
	}
	denominator = math.Sqrt(squaredOne) * math.Sqrt(squaredTwo)
	similarity = numerator / denominator
	return similarity
}

func getAverageRating(totalUsers[]int, movie int, filename string) int {
	var sum float64
	var num int
	testData := getTestData(filename)
	
	for _, value := range totalUsers {
		sum += float64(testData[value][movie]) //sum of all the users
		num++
	}
	return int(sum / float64(num))
}

func getPrediction(activeUser int, activeMovie int, filename string) int {
	trainingData := getData() //get training data for similarity prediction
	kSimilarItems := [15]int{} //array to store the top 15 most similar items
	kSimilarRatings := [15]float64{} //array to store similarity scores for the top 15 most similar items
	leastSimilar := 0 //int variable to keep track of the index for the similar item in kSimilarItems with the smallest similarity score
	
	//find k most similar movies to active movie
	for other := 0; other < 1000; other++ {
		if trainingData[activeUser][other] != 0 { //other movie also has a rating
			similarityScore := itemCosine(activeMovie, other, filename) //returns similarity score between active user and other user

			if similarityScore > kSimilarRatings[leastSimilar] {
				kSimilarItems[leastSimilar] = other //update most similar users array
				kSimilarRatings[leastSimilar] = similarityScore //update most similar users' ratings array
			}
			//traverse the similar items and ratings arrays to find the least similar item
			for i := 0; i < 15; i++ {
				if kSimilarRatings[i] < kSimilarRatings[leastSimilar] {
					leastSimilar = i
				}
			}
		}
	}
	
	//predict ratings for active movie given most similar 15 movies
	var numerator float64
	var denominator float64
	var prediction float64
	
	for u := 0; u < 15; u++ { //traverse through the top 15 most similar movies
		numerator = numerator + (kSimilarRatings[u] * float64(trainingData[activeUser][kSimilarItems[u]]))
		denominator = denominator + kSimilarRatings[u]
	}

	prediction = numerator / denominator
	return int(math.Round(prediction))
} 

//Function implements item-based collaborative filtering algorithm with adjusted cosine similarity
func itemBasedCosinePrediction(filename string)  [100][1000]int {
	testData := getTestData(filename)
	updatedTestData := testData
	counter := 0

	for user := 3; user < 100; user++ {
		for movie := 0; movie < 1000; movie++ {
			if testData[user][movie] == 0 { //need prediction for active user
				prediction := getPrediction(user, movie, filename) //to get prediction for movie item for active user
				fmt.Println(prediction)
				updatedTestData[user][movie] = prediction
			}
		}
		counter++
		fmt.Println("Prediction done for user:", counter)
	}
	return updatedTestData
}