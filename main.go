package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mongodb-go-tutorial/app"
)

func test(str ...string) {
	for i, s := range str {
		fmt.Println(i, s)
	}

}

func main() {

	test("aaa", "bbb", "ccc")

	app.Setup()

	//app.CreateBook(app.Book{ "How to Use Go With MongoDB", "GeeksforGeeks", 100 }) // TODO: 622fd2246de12e1b7f36c4db
	//app.CreateBook(app.Book{ "How to Do CRUD Transactions in MongoDB with Go", "hackajob Staff", 99 }) // TODO: 622fd2246de12e1b7f36c4dc
	//app.CreateBook(app.Book{ "MongoDB Go Driver туториал", "pocoZ", 101 }) // TODO: 622fd2246de12e1b7f36c4dd

	result, err := app.GetBook("622fd2246de12e1b7f36c4dc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name='%v'; Author='%v'; PageCount='%v'; \n", result.Name, result.Author, result.PageCount)

	results, err := app.GetBooks()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(results)

	/*
	* https://www.geeksforgeeks.org/loops-in-go-language
	 */
	//allRec := 0
	for count, result := range results {
		//allRec = count
		//fmt.Printf("Name='%v'; Author='%v'; PageCount='%v'; \n", result.Name, result.Author, result.PageCount)
		fmt.Printf("%v. Name='%v'; Author='%v'; PageCount='%v'; \n", count, result.Name, result.Author, result.PageCount)
	}
	//allRec = allRec * allRec

	/*
	* @see https://serveanswer.com/questions/how-to-convert-string-to-primitive-objectid-in-golang
	*      https://stackoverflow.com/questions/60864873/primitive-objectid-to-string-in-golang
	 */
	bookId, err := primitive.ObjectIDFromHex("622fd2246de12e1b7f36c4db")
	app.UpdateBook(bookId, 199)

	app.DeleteBook(bookId)

	fullName := "GeeksforGeeks"
	resultsAll, errAll := app.FindAuthorBooks(fullName)
	if errAll != nil {
		log.Fatal(errAll)
	}
	for count, result := range resultsAll {
		fmt.Printf("%v. FullName='%v'; > Name='%v'; Author='%v'; PageCount='%v'; \n", count, fullName, result.Name, result.Author, result.PageCount)
	}

}
