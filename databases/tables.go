package databases

type Table struct {
	Name        string
	ColumnNames []string
	Columns     map[string]string
}

func GetTables() []Table {

	/*
	*-------------------------------------------------------
	*					TABLES
	*-------------------------------------------------------
	*Create tables in the database by using the given template.
	*Copy and use this format for future tables
	 */

	tables := []Table{
		{
			"users",
			[]string{"id", "aadhaar_id", "name", "dob", "image_link", "created_at", "updated_at"},
			map[string]string{
				"id":         "INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY",
				"aadhaar_id": "INT(12) UNIQUE NOT NULL",
				"name":       "VARCHAR(30) NOT NULL",
				"dob":        "DATE",
				"image_link": "VARCHAR(65)",
				"created_at": "DATETIME",
				"updated_at": "DATETIME",
			},
		},
	}
	return tables
}
