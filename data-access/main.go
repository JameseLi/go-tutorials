package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

type Album struct {
	Id     int64
	Title  string
	Artist string
	Price  float32
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumByArtist: %q: %v", name, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("closing rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("albumByArtist: %q: %v", name, err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumByArtist: %q: %v", name, err)
	}
	return albums, nil
}

func albumById(id int64) (Album, error) {
	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return album, fmt.Errorf("albumsById %d: no such album", id)
		}
		return album, fmt.Errorf("albumById: %d: %v", id, err)
	}
	return album, nil
}

func addAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func main() {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "172.16.78.12:3306",
		DBName:               "zdf",
		AllowNativePasswords: true,
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Successfully connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	album, err := albumById(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

	albumId, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albumId)
}
