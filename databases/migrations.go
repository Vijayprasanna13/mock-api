package databases

import (
	"log"
	"fmt"
)


func CreateTables() {

	tables := GetTables()
	
	for _, table := range tables{
		var column_desc string
		for  _, column_name := range table.ColumnNames{
			column_desc += column_name+" "+table.Columns[column_name]+","
		} 
		_, err := DB_CONN.Exec("create table if not exists "+table.Name+"("+column_desc[:len(column_desc)-1]+")")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("<"+table.Name+"> table is created.")
	}
	
	fmt.Println("Migration is complete...")
	defer DB_CONN.Close()
}
