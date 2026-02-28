package middleware
import (
	"log"
)

func PanicButton(err error) {
	 if err != nil {
		log.Fatal(err)
	}
}