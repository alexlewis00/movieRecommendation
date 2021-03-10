//Declaring main package, groups functions and all files in same directory movieRecommendation
package main

//importing packages
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
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

	data, err := os.Open("../Data/test5.txt") //read contents of file txt into data array
	if err != nil {
		fmt.Println("Failed to read file")
	}
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)
	//buffer := [5]int{} //make slice for buffer

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	//I can print the strings but I need to be able to access it has a multidimensional array?
	//Need to somehow access the ratings related to the specific userid and movieid
	
}
