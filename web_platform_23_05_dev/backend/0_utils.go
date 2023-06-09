package backend

//

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ErrorHandler(err error) {
	fmt.Printf("error: %s", err.Error())
}

func ErrHandlerWithContext(err error, context *gin.Context) {
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	ErrorHandler(err)
}

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////

//
