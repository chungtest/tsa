package main

import (
  //"net/http"
  "github.com/gin-gonic/gin"
  "errors"
  "strings"
  "database/sql"
  //"fmt"
  "log"
  _"github.com/denisenkom/go-mssqldb"
)

// This isued for the intiial testing of databases as I not sure how to map
// phone numbers correctly
type Contact struct {
  FullName string `json:"full_name"`
  Email string `json:"email"`
  PhoneNumber string `json:"phone_numbers"`
  //Should be the below
  //PhoneNumber string `json:"phone_numbers"`
}

// Final structure should be closer to the below
// Used for the inserting
type ContactNums struct {
  FullName string `json:"full_name"`
  Email string `json:"email"`
  PhoneNumbers []string `json:"phone_numbers"`
}

// The main api server
func main() {
  r := gin.Default()

  // routing could be anything not jusrt /api could be /tsa/api
  // versioning seems to the way to go
  v1 := r.Group("api/v1")
    {
        v1.GET("/contacts", GetContacts)
        v1.POST("/newcontact", NewContact)
        //v1.POST("/updatecontact", UpdateContact)
    }
    r.Run(":8080")
}

func NewContact (c *gin.Context) {

  var input ContactNums

  err := c.Bind(&input)
  if err != nil {
    // likely a malformed request in the JSON
    c.JSON(400, gin.H{"Failed" : err})
  } else {
    // Validate the incoming data. Not sure the best way to deal with that
    messages, err := isValid(input)
    if err != nil {
      c.JSON(422, gin.H{"Failed" : messages})
      return
    }

    //Check if the person exists
    if contactExists(input.FullName) {
      c.JSON(422, gin.H{"Failed" : "Person exists use Update"})
      return
    }

    // sql connection
    db, err := sql.Open("sqlserver", "sqlserver://tsauser:tsapassword@localhost?database=tsa&connection+timeout=30")
  	if err != nil {
  		log.Fatal(err)
    }

    // Is this how you do transactions??
    // If the contact passes but a phone number fails we should not insert and roll everything back
    tx, err := db.Begin()
    if err != nil {
    	log.Fatal(err)
    }
    defer tx.Rollback()
    // user parameterise sql in case of injection
    stmt, err :=  tx.Prepare("insert into tsa.dbo.contact (full_name, email) values (@p1, @p2)")
    if err != nil {
      log.Fatal(err)
    }
    _, err = stmt.Exec(input.FullName, input.Email)
    if err != nil {
      log.Fatal(err)
    } else {
      for _, number := range input.PhoneNumbers {
        stmt, err :=  tx.Prepare("insert into tsa.dbo.contact_phone_number (full_name, number) values (@p1, @p2)")
        if err != nil {
          log.Fatal(err)
        }
        _, err = stmt.Exec(input.FullName, number)
        if err != nil {
          log.Fatal(err)
        }
      }
    }
    err = tx.Commit()
    if err != nil {
    	log.Fatal(err)
    } else {
      c.JSON(200, gin.H{"Success" : "ok"})
    }
  }
}

func isValid (input ContactNums) ([]string, error) {
  var messages []string

  // BIG assumption that there is a space between first and family names
  if strings.Count(input.FullName, " ") != 1 {
    messages = append (messages, "The full_name is must be in (<firstname> <surname>) format")
  }
  if input.FullName == "" {
    messages = append (messages, "The full_name is empty")
  }
  if input.Email == "" {
    messages = append (messages, "The email is empty")
  }
  if input.PhoneNumbers == nil {
    messages = append (messages, "The are no phone numbers")
  }
  if len(messages) > 0 {
    return messages, errors.New("There were empty fields")
  } else {
    return messages, nil
  }
}

// check to see if the person exists
func contactExists (name string) bool {
  db, err := sql.Open("sqlserver", "sqlserver://tsauser:tsapassword@localhost?database=tsa&connection+timeout=30")
	if err != nil {
		log.Fatal(err)
  }

  rows, err := db.Query("select	full_name from tsa.dbo.contact where full_name = @p1", name)
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()
  // I not sure how to check if rows contains anything except for looping
  // This is kind of like the DBReader in c#????
  for rows.Next() {
    return true
  }
  return false
}

// Get data from the database
func GetContacts(c *gin.Context) {
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
 defer rows_withNumbers.Close()
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
  c.JSON(200, output)
  //return output
}
