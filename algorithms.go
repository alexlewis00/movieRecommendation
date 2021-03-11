//Declaring main package, groups functions and all files in same directory movieRecommendation
package main

//importing packages
import (
	"bufio"
	"fmt"
	"io/ioutil"
	//"reflect"
	"strconv"
	"strings"
	//"math"
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
	users := [101][1001]int{}
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
func userCosine(filename string) int {
	//"Similar users rate similarly", need test data to find the similarity of the user to the users of the train data
	//and be able to use that similarity to predict the unknown ratings of the test data users
	trainingData := getData() //Get the test data
	testData := getTestData(filename) //Get the train data
	var totalNum int
	var totalDen int
	var totalNumSum int
	var totalDenSum int

	var similarity int
	
	//traverse through the test data users, each user is the active user the system must recommend movies for
	for users := 1; users < 101; users++ {
		for movies := 1; movies < 1001; movies++ {
			for testData[users][movies] != 0 { //for the active users of the test data that have movie ratings
			
			}
		}
	}
	for activeUser := 1; activeUser < 101; activeUser++ { //for active user 1 to 100
		for activeMovie := 1; activeMovie < 1001; activeMovie++ {
			if testData[activeUser][activeMovie] != 0 {
				for trainUser := 1; trainUser < 201; trainUser++ {
					if trainingData[trainUser][activeMovie] != 0 { //this means we are dealing with an active user and train user that both have rated the same movie
						//for the numerator
						product := 0
						product = testData[activeUser][activeMovie] * trainingData[trainUser][activeMovie]
						totalNum += product
						product = (testData[activeUser][activeMovie]/2) * (trainingData[activeUser][activeMovie]/2)
						totalDen += product
					}
					totalNumSum += totalNum
					totalDenSum += totalDen
				}
			}
		}
	}
	return similarity
}
