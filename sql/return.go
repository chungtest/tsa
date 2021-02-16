package main

import (
  "database/sql"
  "fmt"
  "log"
  _"github.com/denisenkom/go-mssqldb"
)

type Contact struct {
  FullName string `json:"full_name" gorm:"primary_key"`
  Email string `json:"email"`
  PhoneNumber string `json:"phone_numbers"`
  //Should be the below
  //PhoneNumber string `json:"phone_numbers"`
}

func main() {
  blah := ReturnContacts()
  fmt.Println(blah)
}

// functions?
func ReturnContacts() []Contact {
    //connect to the database
  db, err := sql.Open("sqlserver", "sqlserver://tsauser:tsapassword@localhost?database=tsa&connection+timeout=30")
	if err != nil {
		log.Fatal(err)
  }

  //execute a query for testing
 rows_withNumbers, err := db.Query(`
    select	con.full_name, con.email, ph.number
    from tsa.dbo.contact con
	     inner join tsa.dbo.contact_phone_number ph on ph.full_name = con.full_name
    order by con.full_name`)
 if err != nil {
   log.Fatal(err)
 }

 var output []Contact
 //var previousContact string = ""

 // return rows, expecting two at this moment
 for rows_withNumbers.Next() {
   var full_name string
   var email string
   var phone_number string
   err := rows_withNumbers.Scan(&full_name, &email, &phone_number)
   if err != nil {
     log.Fatal(err)
   }

   // mapping it to the struct
   contact := Contact{full_name, email, phone_number}
   output = append(output, contact)
   // Manually map to the struct
   // Need to reduce the duplicate
   //if full_name != previousContact {
     // Map to the a new contact as the person doesn't exist
  ///} else {
     // map to exisitng person and only attach the phone number
   //}

 }

  defer db.Close()
  return output
}
