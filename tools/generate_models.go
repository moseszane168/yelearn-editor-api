package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "Database host")
	port := flag.Int("port", 5433, "Database port")
	user := flag.String("user", "postgres", "Database user")
	password := flag.String("password", "password", "Database password")
	dbname := flag.String("dbname", "editor", "Database name")
	tableNames := flag.String("tables", "", "Comma-separated list of table names to generate models for")
	flag.Parse()

	if *dbname == "" {
		log.Fatal("Please provide a database name using the -dbname flag")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*host, *port, *user, *password, *dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var tables []string
	if *tableNames == "" {
		tables, err = getAllTables(db)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tables = strings.Split(*tableNames, ",")
	}

	dir := fmt.Sprintf("model/%s", *dbname)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {
		filePath := filepath.Join(dir, fmt.Sprintf("%s.go", table))
		modelsFile, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create file for table %s: %v", table, err)
		}
		defer modelsFile.Close()

		fmt.Fprintln(modelsFile, "package model")
		fmt.Fprintln(modelsFile, "")
		fmt.Fprintln(modelsFile, "import (\"time\")")

		err = generateModel(db, modelsFile, table)
		if err != nil {
			log.Fatalf("Failed to generate model for table %s: %v", table, err)
		}

		fmt.Printf("Model for table %s generated successfully in %s\n", table, filePath)
	}
}

func getAllTables(db *sql.DB) ([]string, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE' AND table_schema = 'public'
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, rows.Err()
}

func generateModel(db *sql.DB, file *os.File, table string) error {
	query := fmt.Sprintf(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = '%s'
	`, table)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Fprintf(file, "\n// %s represents the %s table\ntype %s struct {\n", CamelCase(table), table, CamelCase(table))
	for rows.Next() {
		var columnName, dataType string
		err = rows.Scan(&columnName, &dataType)
		if err != nil {
			return err
		}
		fmt.Fprintf(file, "    %s %s `json:\"%s\" db:\"%s\"`\n", CamelCase(columnName), GoType(dataType), columnName, columnName)
	}
	fmt.Fprintln(file, "}")

	return rows.Err()
}

// CamelCase converts a string to camel case.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

// GoType maps PostgreSQL data types to Go data types.
func GoType(pgType string) string {
	switch pgType {
	case "integer", "smallint", "bigint":
		return "int"
	case "numeric", "real", "double precision":
		return "float64"
	case "boolean":
		return "bool"
	case "character", "character varying", "text":
		return "string"
	case "timestamp", "timestamp without time zone", "timestamp with time zone", "date":
		return "time.Time"
	default:
		return "interface{}"
	}
}
