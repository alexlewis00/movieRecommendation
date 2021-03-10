//Declaring main package, groups functions and all files in same directory movieRecommendation
package main

//importing packages
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

//run file in terminal with "go run <filename>" command
func main() {
	//Step 1: Read training data and insert into a 2d array
	trainingData := getData()
	fmt.Println(trainingData[0][5]) //prints userid 0's rating on movieid 5
	userCosine()
}

//Function to read the training data and insert data into 2d array (200 users x 1000 movies)
func getData() [200][1000]int { //func <function name> <returning value of specified type: 2d array of integers>
	data, err := ioutil.ReadFile("../Data/train.txt") //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read file")
		return [200][1000]int{}
	}
	//scan input as sequence of space-delimited tokens
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Split(bufio.ScanWords) //split function for only scanning words, not spaces
	buffer := [200][1000]int{} //buffer to store data in

	for row := 0; row < 200; row++ {
		for col := 0; col < 1000; col++ {
			scanner.Scan() //advances scanner to next token
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

/*
Task: 
Design and develop collaborative filtering algorithms that predict the unknown ratings in the test data
by learning users' preferences from the training data
*/

//Function for the user-based collaborative filtering algorithm with Cosine Similarity
func userCosine() {
	/* Process:
		1. Consider active user a from test data
		2. Find k other users from training data who have similar ratings to a's ratings
			- Use Cosine Similarity to determine similarity between active user a and other users from training data
			- Sort by highest order of most similar users (high rating similarity)
			- Choose some k number of similar users for predicition calculation of user a's other ratings
		3. Estimate a's other ratings based on the ratings of the k similar users
	*/

	/* Step 1: Consider active user a from test data
		1. Need to read and scan test data
		2. Iterate through data and compare [row][column] pairs with testing data, userid and movieid
	*/
	//create another 2d array of test data
	data, err := ioutil.ReadFile("../Data/test5.txt") //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read file")
	}
	//scan input as sequence of space-delimited tokens
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Split(bufio.ScanWords) //split function for only scanning words, not spaces
	for scanner.Scan() {
		fmt.Println(scanner.Text()) //printing string every line of data
	}
	/*
	Need to compare active user a being the test data of the userid with the movies that that user has rate
	Compare this with the training data of that specific movie id
	Find similarity between ratings in order to find similar users
	Just need a way of extracting that rating with the movie id pairs
	Creating an array of arrays? [7687][3] [userids][3] -> [3] = [201, 237, 4]? Need to be able to access each number individually
	*/
}
