package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mongodb-go-tutorial/app"
)

func main() {
	app.Setup()

	//author := Author{ "GeeksforGeeks" }

	//newBook1 := Book{ "How to Use Go With MongoDB", "GeeksforGeeks", 100 } // TODO: 622f6e3fc134f9ed331b4b33
	//newBook2 := Book{ "How to Do CRUD Transactions in MongoDB with Go", "hackajob Staff", 99 } // TODO: 622f6e3fc134f9ed331b4b34
	//newBook3 := Book{ "MongoDB Go Driver туториал", "pocoZ", 101 } // TODO: 622f6e3fc134f9ed331b4b35
	//CreateBook(newBook1)
	//CreateBook(newBook2)
	//CreateBook(newBook3)

	//result, err := GetBook("622f6e3fc134f9ed331b4b33")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Name='%v'; Author='%v'; PageCount='%v'; \n", result.Name, result.Author, result.PageCount)

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
	bookId, err := primitive.ObjectIDFromHex("622f6e3fc134f9ed331b4b34")
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
