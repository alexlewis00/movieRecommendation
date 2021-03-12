//Declaring main package, groups functions and all files in same directory movieRecommendation
package main

//importing packages
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"sort"
	//"reflect"
	"math"
	"strconv"
	"strings"
)

//run file in terminal with "go run <filename>" command
func main() {
	fmt.Println(userBasedPrediction())
}

//Function to read the training data and insert data into 2d array (200 users x 1000 movies)
func getData() [201][1001]int { //func <function name> <returning value of specified type: 2d array of integers>
	data, err := ioutil.ReadFile("../Data/train.txt") //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read training file")
		return [201][1001]int{}
	}
	//scan input as sequence of space-delimited tokens
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Split(bufio.ScanWords) //split function for only scanning words, not spaces
	buffer := [201][1001]int{}     //buffer to store data in

	for row := 1; row < 201; row++ {
		for col := 1; col < 1001; col++ {
			scanner.Scan()                                   //advances scanner to next token
			currentRank, err := strconv.Atoi(scanner.Text()) //Atoi: string conversion to int
			if err != nil {
				fmt.Println("Error with scanner at token")
				return [201][1001]int{}
			}
			buffer[row][col] = currentRank //adding integer value pos into buffer at specified position
		}
	}
	return buffer
}

//Function to read the test data passed as an argument and insert data into 2d array (100 users x 3 attributes)
func getTestData(filename string) [101][1001]int {
	data, err := ioutil.ReadFile(filename) //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read test file")
		return [101][1001]int{}
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data))) //scan input as string of space-delimited tokens
	scanner.Split(bufio.ScanLines)                               //split function for only scanning words, not spaces
	checkUsers := []int{}
	var checker int
	users := [101][1001]int{} //100 users starting at userid 1, 1000 movies starting at movieid 1
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
		users[len(checkUsers)][movieId] = rating
	}
	return users
}

/*
Task:
Design and develop collaborative filtering algorithms that predict the unknown ratings in the test data
by learning users' preferences from the training data
*/

//Function for the user-based collaborative filtering algorithm with Cosine Similarity
func userCosine(activeUser int, filename string) [201]float64 {
	/*
		"Similar users rate similarly", need test data to find the similarity of the active user to the users of the train data via cosine similarity
		and be able to use that similarity to predict the unknown ratings of the test data users by returning the array with the top 20 most similar users
	*/
	trainingData := getData() //Get the test data
	testData := getTestData(filename) //Get the train data, set for the active user
	product := 0
	squaredOne := 0
	squaredTwo := 0
	activeRated := []int{} //keep track of all movies active user has rated
	similarK := [201]float64{} //array to keep similarity scores of other similar users where index is that user's id

	//traverse through movies of active user to see the movies this user has actually rated
	for activeMovie := 1; activeMovie < 1001; activeMovie++ {
		if testData[activeUser][activeMovie] != 0 { //that means there is a rating
			activeRated = append(activeRated, activeMovie)
		}
	}
	//traverse through the training data to find other users who have rated at least one of the same movies as the active user in order to perform cosine similarity
	for otherUser := 1; otherUser < 201; otherUser++ {
		checker := 0 //check if other user has rated at least one of the same movies
		for _, movieID := range activeRated { //iterates over activeRated slice which is the movies the active user has rated, movieID stores the specific value at the index of the iteration
			if trainingData[otherUser][movieID] != 0 { //other user has rated that movie
				checker = 1
				product = product + (testData[activeUser][movieID] * trainingData[otherUser][movieID])
				squaredOne = squaredOne + (testData[activeUser][movieID] * testData[activeUser][movieID])
				squaredTwo = squaredTwo + (trainingData[otherUser][movieID] * trainingData[otherUser][movieID])
			}
		}
		if checker > 0 { //if there is a similarity to calculate
			numerator := float64(product)
			denominator := math.Sqrt(float64(squaredOne)) * math.Sqrt(float64(squaredTwo))
			denominator = math.Round(denominator/0.05) * 0.05 //rounds denominator float value to nearest 2 decimal points
			similarity := float64(numerator / denominator) //similarity of active user with other user of index in for loop iteration
			similarity = math.Round(similarity/0.0005) * 0.0005
			similarK[otherUser] = similarity //records the similarity score of the other user at the index of the user's id
		}
	}
	//sort top 20 most similar users to active user based on similarity rating
	ratings := []float64{} //temp slice to store similarity ratings while sorting top 20 most similar users
	for i := 0; i < 201; i++ { //traverse through similarK
		if similarK[i] != 0 { //similarity score present for that other user
			ratings = append(ratings, similarK[i])
		}
	}
	sort.Float64s(ratings)  //sorts similarity scores in increasing order
	ratingsK := []float64{} //slice to store the top 20 similarity scores
	//traverse through the ratings' slice to get the top 20 similarity scores
	for j := len(ratings) - 1; j > len(ratings)-21; j-- {
		ratingsK = append(ratingsK, ratings[j]) //append the top 20 similarities into the ratingsK
	}
	for z := 1; z < 201; z++ { //traverse through the similarity ratings of similarK
		check := 0
		if similarK[z] != 0 { //there is a similarity score
			for index := range ratingsK {
				if similarK[z] == ratingsK[index] { //if the similarity rating equals one of the top 20 ratings
					check = 1
				}
			}
			//if the similarity rating is not one of the top 20 ratings, if it is not found in ratingsK
			if check == 0 {
				similarK[z] = 0 //set rating to 0, similarity rating not one of the top 20, user won't be applied into prediction
			}
		}
	}
	return similarK //return similarity scores of the top 20 most similar users to the given active user
}

//Function for the user-based collaborative filtering algorithm with Pearson Correlation

func userBasedPrediction() [101][1001]int {
	//Using Cosine Similarity
	test5 := "../Data/test5.txt"
	testData := getTestData("../Data/test5.txt")
	trainingData := getData()
	var numerator float64
	var denominator float64
	var prediction int

	for activeUser := 1; activeUser < 101; activeUser++ { //traverse through the test data, set of active users we predict for
		for activeMovie := 1; activeMovie < 1001; activeMovie++ {
			check := 0
			if testData[activeUser][activeMovie] == 0 { //need to make prediction for the active user a
				check = 1 //prediction needed
				kUsers := userCosine(activeUser, test5) //userCosine function returns an array of 200 users with only similarity ratings of the top 20 most similar users based on cosine similarity
				for i := 1; i < 201; i++ {
					if kUsers[i] != 0 { //meaning that there is a similarity rating at that index, the userid has a similarity rating with the active user
						numerator = numerator + (kUsers[i] * float64(trainingData[i][activeMovie]))
						numerator = math.Round(numerator/0.0005) * 0.0005
						denominator = denominator + kUsers[i]
						denominator = math.Round(denominator/0.0005) * 0.0005
					}
				}
			}
			if check > 0 { //execute prediction
				prediction = int(numerator/denominator) //float64 to int value rounds towards zero
				testData[activeUser][activeMovie] = int(prediction) //update predicted rating in the testData for the active user's active movie
			}
		}
	}
	//Using Pearson Correlation
	return testData
}
