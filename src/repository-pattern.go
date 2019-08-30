package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	sq "github.com/Masterminds/squirrel"
)

// Employee ...
type Employee struct {
	ID   int
	Name string
	City string
}

// NewDB ...
func NewDB() (*sql.DB, error) {
	driver := "mysql"
	user := "withgoods"
	pass := "withgoods"
	database := "goblog"

	db, err := sql.Open(driver, user+":"+pass+"@/"+database)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// EmployeeRepository ...
type IEmployeeRepository interface {
	GetEmployees() ([]*Employee, error)
}

// BaseRepository ...
type Repository struct {
	db *sql.DB
}

type IEmployeeService interface {
	GetEmployeesExcept(names []string) ([]*Employee, error)
}

type EmployeeService struct {
	repo IEmployeeRepository
}

func (service *EmployeeService) GetEmployeesExcept(names []string) ([]*Employee, error) {
	employees, err := service.repo.GetEmployees()

	if err != nil {
		return nil, err
	}

	filtered := make([]*Employee, 0)

	for _, employee := range employees {
		for _, name := range names {
			if name != employee.Name {
				filtered = append(filtered, employee)
			}
		}
	}

	return filtered, nil
}

// GetEmployees ...
func (repo *Repository) GetEmployees() ([]*Employee, error) {
	defer repo.db.Close()

	query, _, err := sq.Select("id, name, city").From("employee").OrderBy("id desc").ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, err
	}

	employees := make([]*Employee, 0)

	for rows.Next() {
		var id int
		var name, city string
		err = rows.Scan(&id, &name, &city)
		employee := Employee{
			ID:   id,
			Name: name,
			City: city,
		}
		employees = append(employees, &employee)
	}

	return employees, nil
}
