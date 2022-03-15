package main

import (
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"mongodb-go-tutorial/app"
	"mongodb-go-tutorial/models"
	"strconv"
	"time"
)

func test(str ...string) {
	for i, s := range str {
		fmt.Println(i, s)
	}

}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

/**
 * @see https://stackoverflow.com/questions/24987131/how-to-parse-unix-timestamp-to-time-time
 *      https://golang.cafe/blog/golang-time-format-example.html
 */
func ToString(date primitive.Timestamp, layout string) string {
	uintDate := strconv.FormatUint(uint64(date.T), 10)
	intDate, err := strconv.ParseInt(uintDate, 10, 64)

	if err != nil {
		panic(err)
	}
	unixDate := time.Unix(intDate, 0)

	return unixDate.Format(layout)
}

func Print(user models.User) {
	ObjectID := user.ObjectID.Hex()
	Username := user.Username
	Email := user.Email
	Password := user.Password
	CreatedAt := ToString(user.CreatedAt, time.RFC822)
	DeactivatedAt := ToString(user.DeactivatedAt, time.RFC822)
	fmt.Printf("\nObjectID='%v'; Username='%v'; Email='%v'; Password='%v'; CreatedAt='%v'; DeactivatedAt='%v'; \n\n",
		ObjectID, Username, Email, Password, CreatedAt, DeactivatedAt)
}

func PrintList(u int, user models.User) {
	ObjectID := user.ObjectID.Hex()
	Username := user.Username
	Email := user.Email
	Password := user.Password
	CreatedAt := ToString(user.CreatedAt, time.RFC822)
	DeactivatedAt := ToString(user.DeactivatedAt, time.RFC822)
	fmt.Printf("%v. ObjectID='%v'; Username='%v'; Email='%v'; Password='%v'; CreatedAt='%v'; DeactivatedAt='%v'; \n",
		u, ObjectID, Username, Email, Password, CreatedAt, DeactivatedAt)
}

func main() {

	test("aaa", "bbb", "ccc")

	rand.Seed(time.Now().UnixNano())
	s := RandomString(22)
	se := base64.StdEncoding.EncodeToString([]byte(s))
	log.Printf("New Password: '%v' == '%v' \n", s, se)

	app.Setup()

	/**
	 * @see https://stackoverflow.com/questions/60864873/primitive-objectid-to-string-in-golang
	 *      https://golangify.com/unix-timestamp
	 *      https://yourbasic.org/golang/convert-string-to-byte-slice
	 *      https://golangdocs.com/generate-random-string-in-golang
	 *      https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	 */
	//
	password := RandomString(22)
	newUser := models.User{
		primitive.NewObjectID(),
		"Alex",
		"alex@smarttrader.com.ua",
		base64.StdEncoding.EncodeToString([]byte(password)),
		primitive.Timestamp{T: uint32(time.Now().Unix())},
		primitive.Timestamp{T: uint32(time.Now().Unix())},
	}

	strNewUserId, err := app.CreateUser(newUser)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New UserID:", strNewUserId)

	user, err := app.GetUser(strNewUserId)
	if err != nil {
		log.Fatal(err)
	}
	Print(user)

	users, err := app.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	for u, user := range users {
		PrintList(u, user)
	}

	//
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
