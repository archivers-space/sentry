package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Url struct {
	Url           *url.URL
	Created       time.Time
	Updated       time.Time
	LastGet       time.Time
	Host          string
	Status        int
	ContentType   string
	ContentLength int64
	Title         string
}

// ShouldFetch returns weather the url should be added to the queue for updating
// should return true if the url is new, or if we haven't checked this url in a while
func (u *Url) ShouldEnqueueGet() bool {
	return (u.LastGet.IsZero() || time.Since(u.LastGet) > cfg.StaleDuration) && !enqued[u.Url.String()]
}

func (u *Url) ShouldEnqueueHead() bool {
	return (u.Created == u.Updated || u.LastGet.IsZero() || time.Since(u.Updated) > cfg.StaleDuration) && !enqued[u.Url.String()]
}

func (u *Url) Read(db sqlQueryable) error {
	if u.Url != nil {
		row := db.QueryRow(fmt.Sprintf("select %s from urls where url = $1", urlCols()), u.Url.String())
		return u.UnmarshalSQL(row)
	}

	return ErrNotFound
}

// func (u *Url) ReadDomain(db sqlQueryable) error {
// 	if u.Url == nil {
// 		return ErrNotFound
// 	}

// 	d := &Domain{
// 		Host: u.Url.Host,
// 	}

// 	if err := d.Read(db); err != nil {
// 		return err
// 	}

// 	u.Host = d
// 	return nil
// }

func (u *Url) Insert(db sqlQueryExecable) error {
	u.Created = time.Now()
	u.Updated = u.Created
	u.Url = NormalizeURL(u.Url)
	u.Host = u.Url.Host
	_, err := db.Exec(fmt.Sprintf("insert into urls (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", urlCols()), u.SQLArgs()...)
	return err
}

func (u *Url) Update(db sqlQueryExecable) error {
	u.Updated = time.Now()
	u.Url = NormalizeURL(u.Url)
	if u.ContentLength < -1 {
		u.ContentLength = -1
	}
	if u.Status < -1 {
		u.Status = -1
	}
	_, err := db.Exec("update urls set created=$2, updated=$3, last_get=$4, host=$5, status=$6, content_type=$7, content_length=$8, title=$9 where url = $1", u.SQLArgs()...)
	return err
}

func (u *Url) Delete(db sqlQueryExecable) error {
	_, err := db.Exec("delete from urls where url = $1", u.Url.String())
	if err != nil {
		logger.Println(err, u)
	}
	return err
}

// DocLinks extracts a page's linked documents
// extracts all a[href] links from a qoquery document.
func (u *Url) DocLinks(doc *goquery.Document) ([]*Link, error) {
	links := make([]*Link, 0)
	// generate a list of normalized links
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")

		// Resolve address to source url
		address, err := u.Url.Parse(val)
		if err != nil {
			logger.Printf("error: resolve URL %s - %s\n", val, err)
			return
		}

		// allocate normalized link
		l := &Link{
			Src: u,
			Dst: &Url{
				Url: NormalizeURL(address),
			},
		}

		links = append(links, l)
	})

	return links, nil
}

func urlCols() string {
	return "url, created, updated, last_get, host, status, content_type, content_length, title"
}

func (u *Url) UnmarshalSQL(row sqlScannable) error {
	var (
		rawurl, host, mime, title         string
		created, updated, lastGet, length int64
		status                            int
	)

	if err := row.Scan(&rawurl, &created, &updated, &lastGet, &host, &status, &mime, &length, &title); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		logger.Println(err)
		return err
	}

	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	*u = Url{
		Created:       time.Unix(created, 0),
		Updated:       time.Unix(updated, 0),
		LastGet:       time.Unix(lastGet, 0),
		Url:           parsedUrl,
		Host:          host,
		Status:        status,
		ContentType:   mime,
		ContentLength: length,
		Title:         title,
	}

	return nil
}

func (u *Url) SQLArgs() []interface{} {
	t := int64(0)
	if !u.LastGet.IsZero() {
		t = u.LastGet.Unix()
	}
	return []interface{}{
		u.Url.String(),
		u.Created.Unix(),
		u.Updated.Unix(),
		t,
		u.Host,
		u.Status,
		u.ContentType,
		u.ContentLength,
		u.Title,
	}
}
