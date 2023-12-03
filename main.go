package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Language struct {
	ID       int64
	Host     string
	Target   string
	Priority float32
}

type FlashCard struct {
	ID         int64
	front      string
	back       string
	hint       string
	difficulty int64
	lang       Language
}

type Album struct {
	title  int64
	artist string
	price  string
}

func main() {
	// Capture connection properties.
	//OS property setting on windows is scuffed; learn better
	cfg := mysql.Config{
		User:   "root",
		Passwd: "nunya!",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	//Init ping
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	languages, err := allLanguages()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Languages found: %v\n", languages)

	//Here we can decide to load all flashcards at once or to put them in different calls
	//As this is imple data, we can probably get all cards and store in local memory. But to utilize GraphQL/ etc, should move on
	reviewCards, err := allCards()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Cards found: %v\n", reviewCards)

	// //Run the method call with the Album Struct
	// albums, err := albumsByArtist("John Coltrane")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Albums found: %v\n", albums)
	// // Hard-code ID 2 here to test the query.
	// alb, err := albumByID(2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Album found: %v\n", alb)

}

func allLanguages() ([]Language, error) {
	var languages []Language

	rows, err := db.Query("SELECT * FROM languages")
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var lang Language
		if err := rows.Scan(&lang.ID, &lang.Host, &lang.Target, &lang.Priority); err != nil {
			return nil, fmt.Errorf("albumsByArtist: %v", err)
		}
		languages = append(languages, lang)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	return languages, nil
}

func allCards() ([]FlashCard, error) {
	var cards []FlashCard

	rows, err := db.Query("SELECT * FROM flashcards")
	if err != nil {
		return nil, fmt.Errorf("flashcard: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var card FlashCard
		if err := rows.Scan(&card.ID, &card.front, &card.back, &card.hint, &card.difficulty, &card.lang); err != nil {
			return nil, fmt.Errorf("allCards: %v", err)
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("allCards: %v", err)
	}
	return cards, nil
}

// albumsByArtist queries for albums that have the specified artist name.
// func albumsByArtist(name string) ([]Album, error) {
// 	// An albums slice to hold data from returned rows.
// 	var albums []Album

// 	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
// 	if err != nil {
// 		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
// 	}
// 	defer rows.Close()
// 	// Loop through rows, using Scan to assign column data to struct fields.
// 	for rows.Next() {
// 		var alb Album
// 		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
// 			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
// 		}
// 		albums = append(albums, alb)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
// 	}
// 	return albums, nil
// }

//	adds the specified flashcard to the database,
//
// returning the album ID of the new entry
func addFlashcard(card FlashCard) (int64, error) {
	//TODO: Slim down this or pursue a more persistant style rep
	result, err := db.Exec("INSERT INTO flashcards (front, back, hint, lang) VALUES (?, ?, ?, ?)", card.front, card.back, card.hint, card.lang)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
