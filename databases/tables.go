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
			[]string{"id", "rollno", "name", "password", "department", "is_tronix", "is_appdev", "is_webdev", "is_algos", "is_admin", "contact", "facebook_id", "created_at", "updated_at"},
			map[string]string{
				"id":          "INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY",
				"rollno":      "INT(9)",
				"name":        "VARCHAR(30)",
				"password":    "VARCHAR(70)",
				"department":  "ENUM('chem', 'civil', 'cse', 'eee', 'ece', 'ice', 'mech', 'meta', 'prod')",
				"is_tronix":   "BOOLEAN",
				"is_appdev":   "BOOLEAN",
				"is_webdev":   "BOOLEAN",
				"is_algos":    "BOOLEAN",
				"is_admin":    "BOOLEAN",
				"contact":     "VARCHAR(10)",
				"facebook_id": "VARCHAR(100)",
				"created_at":  "DATETIME",
				"updated_at":  "DATETIME",
			},
		},
	}
	return tables
}
