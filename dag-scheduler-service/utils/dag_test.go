package utils

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var filename string

func init(){
	flag.StringVar(&filename, "f", "", "")
}

func TestCreateDAG(t *testing.T){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dag, err := CreateDAG(ctx, filename)
	assert.Nil(t, err)

	dagJson, _ := json.MarshalIndent(dag, "", " ")
	fmt.Printf("dag = %v\n", string(dagJson))
}
