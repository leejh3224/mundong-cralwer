package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gocolly/colly"
)

// Book ...
type Book struct {
	Title       string
	Author      string
	Category    string
	Publisher   string
	PublishedAt string
	Price       uint32
	Img         string
}

func main() {
	c := colly.NewCollector()

	books := make([]Book, 0, 50)

	c.OnHTML(".list_type_01 > li", func(e *colly.HTMLElement) {
		price := e.ChildText(".info dd:nth-of-type(6)")
		price = strings.Replace(price[:utf8.RuneCountInString(price)-1], ",", "", -1)
		numPrice, _ := strconv.Atoi(price)

		book := Book{
			Title:       e.ChildText(".info dd:nth-of-type(1)"),
			Author:      e.ChildText(".info dd:nth-of-type(2)"),
			Category:    e.ChildText(".info dd:nth-of-type(3)"),
			Publisher:   e.ChildText(".info dd:nth-of-type(4)"),
			PublishedAt: e.ChildText(".info dd:nth-of-type(5)"),
			Price:       uint32(numPrice),
			Img:         e.ChildAttr("img", "src"),
		}

		books = append(books, book)
	})

	for page := 0; page < 3; page++ {
		link := fmt.Sprintf("https://www.munhak.com/book/?dtype=new&ltype=1&page=%d", page)
		c.Visit(link)
	}

	for _, book := range books {
		fmt.Println(book)
	}

	db, err := NewDB()

	if err != nil {
		panic(err.Error())
	}

	var employeeRepo IEmployeeRepository = &Repository{db: db}
	var employeeService IEmployeeService = &EmployeeService{repo: employeeRepo}

	employees, err := employeeService.GetEmployeesExcept([]string{"john"})

	if err != nil {
		panic(err.Error())
	}

	for _, employee := range employees {
		fmt.Println(employee)
	}
}
