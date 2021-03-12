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
	//Step 1: Read training data and insert into a 2d array
	trainingData := getData() //userid: 1-200 , movieid: 1-1000
	testData5 := getTestData("../Data/test5.txt") //userid starts at 1 , movieid starts at 1
	//testData10 := getTestData("../Data/test10.txt")
	//testData20 := getTestData("../Data/test20.txt")
	fmt.Println(trainingData[1][5]) //prints userid 1's rating on movieid 5
	fmt.Println(testData5[1][111])
	//userCosine("../Data/test5.txt")
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
	buffer := [201][1001]int{} //buffer to store data in

	for row := 1; row < 201; row++ {
		for col := 1; col < 1001; col++ {
			scanner.Scan() //advances scanner to next token
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
	scanner.Split(bufio.ScanLines) //split function for only scanning words, not spaces
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

/* Process:
	1. Consider active user a
	2. Find k other users whose ratings are "similar" to a's ratings
		- Use Cosine Similarity to determine similarity between active user a and other user
		- Sort by highest order of most similar users (high rating similarity)
		- Choose some k number to use for rating prediction
	3. Estimate a's ratings based on the ratings of the k similar users
*/

//Function for the user-based collaborative filtering algorithm with Cosine Similarity
func userCosine(activeUser int, filename string) [201]int {
	//"Similar users rate similarly", need test data to find the similarity of the user to the users of the train data
	//and be able to use that similarity to predict the unknown ratings of the test data users
	trainingData := getData() //Get the test data
	testData := getTestData(filename) //Get the train data, activeUser received applies here
	product := 0
	squaredOne := 0
	squaredTwo := 0
	numerator := 0
	denominator := 0
	similarity := 0 //similarity score
	activeRated := []int{} //keep track of all movies active user has rated
	similarK := [201]int{} //array to keep similarity scores of other users where index is the user's id

	//traverse through movies of active user to see the movies this user has actually rated
	for activeMovie := 1; activeMovie < 1001; activeMovie++ { //traverse through the movies for the active user
		//see what movies active user has ranked so can determine similarity to other users
		if testData[activeUser][activeMovie] != 0 { //that means there is a rating
			activeRated = append(activeRated, activeMovie)
		}	
	}
	
	//traverse through the training data to find other users who have rated the same movies in order to perform cosine similarity
	for otherUser := 1; otherUser < 201; otherUser++ {
		checker := 0;
		for _, movieID := range activeRated {
			if trainingData[otherUser][movieID] != 0 { //other user has rated that movie as well
				checker = 1;
				product = product + (testData[activeUser][movieID] * trainingData[otherUser][movieID])
				squaredOne = squaredOne + (testData[activeUser][movieID] * testData[activeUser][movieID])
				squaredTwo = squaredTwo + (trainingData[otherUser][movieID] * trainingData[otherUser][movieID])
			}
		}
		if checker > 0 { //if there is a new similarity
			numerator = product
			denominator = int(math.Sqrt(float64(squaredOne)) * (math.Sqrt(float64(squaredTwo))))
			similarity = numerator / denominator //similarity of active user with other user of index in for loop iteration
			similarK[otherUser] = similarity //records the similarity score of the other user at the index of the user's id
		}
	}
	ratings := []int{}
	for i := 0; i < 201; i++ { //traverse through similarK, put all ratings into a slice in order to order
		ratings = append(ratings, similarK[i])
	}
	sort.Ints(ratings) //sorts similarity scores in increasing order
	ratingsK := []int{}
	for j := len(ratings)-1; j > len(ratings)-21; j-- {
		ratingsK = append(ratingsK, ratings[j]) //append the top 20 similarities into the ratingsK
	}
	for z := 1; z < 201; z++ { //traverse through the similarity ratings of similarK
		check := 0
		for index, _ := range ratingsK {
			if similarK[z] == ratingsK[index] {
				check = 1
			}
		}
		//if the similarity rating is not one of the top 20 ratings, if it is not found in ratingsK
		if check == 0 {
			similarK[z] = 0 //set rating to 0, similarity rating won't be applied
		}
	}
	return similarK //return similarity scores of the top 5 most similar users to active user
}

func userBasedPrediction() {
	//Using Cosine Similarity
	test5 := "../Data/test5.txt"
	testData := getTestData("../Data/test5.txt")
	trainingData := getData()
	numerator := 0
	denominator := 0
	var prediction int

	for activeUser := 1; activeUser < 101; activeUser++ {
		activeMovie := 0
		check := 0
		for activeMovie := 1; activeMovie < 1001; activeMovie++ {
			if testData[activeUser][activeMovie] == 0 { //need to make prediction for the active user a
				check = 1
				kUsers := userCosine(activeUser, test5) //returns k users based on cosine similarity -> 1-200 userids with each index resulting in a similarity rating
				for i := 1; i < 201; i++ {
					if kUsers[i] != 0 { //meaning that there is a similarity rating at that index, the userid has a similarity rating with the active user
						numerator = numerator + (kUsers[i] * trainingData[i][activeMovie])
						denominator = denominator + kUsers[i]
					}
				}
			}
		}
		if check > 0 {
			prediction = numerator / denominator
			testData[activeUser][activeMovie] = prediction //update predicted rating for the active user in the testData
		} 
	}

	//Using Pearson Correlation
}
