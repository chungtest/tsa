package main //required of any executable program

import ( //our Go packages for this project
    "database/sql"
    "fmt"
    "log"
    _"github.com/denisenkom/go-mssqldb"
)


// where our program begins
func main() {
    //connect to the database
    db, err := sql.Open("sqlserver", "sqlserver://tsauser:tsapassword@localhost?database=tsa&connection+timeout=30")
	if err != nil {
		log.Fatal(err)
  }

   //execute a query for testing
  rows, err := db.Query("select full_name, email from tsa.dbo.contact")
	if err != nil {
		log.Fatal(err)
	}

  // return rows, expecting two at this moment
	for rows.Next() {
    var full_name string
    var email string
		err := rows.Scan(&full_name, &email)
		if err != nil {
			log.Fatal(err)
		}
    fmt.Println(full_name, " email >", email)
	}

  //execute a query for testing
 rows_withNumbers, err := db.Query(`
    select	con.full_name, con.email, ph.number
    from tsa.dbo.contact con
	     inner join tsa.dbo.contact_phone_number ph on ph.full_name = con.full_name`)
 if err != nil {
   log.Fatal(err)
 }

 // return rows, expecting two at this moment
 for rows_withNumbers.Next() {
   var full_name string
   var email string
   var phone_number string
   err := rows_withNumbers.Scan(&full_name, &email, &phone_number)
   if err != nil {
     log.Fatal(err)
   }
   fmt.Println("Name > ", full_name, " email >", email, " Phone Number > ", phone_number)
 }

  defer db.Close()

	fmt.Println("Press the Enter Key to terminate the console screen!")
	fmt.Scanln() // wait for Enter Key
}
