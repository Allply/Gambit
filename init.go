package main

import (
  "os"
  "context"
  "encoding/json"
  "fmt"
  "github.com/weaviate/weaviate-go-client/v4/weaviate"
  "github.com/weaviate/weaviate/entities/models"
  "github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)
func main() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	if(len(os.Args) < 2){
		fmt.Print("This is a useless program without an arg\ntry one of the following: \n\n--init      Initialize Schema\n--delete    Delete Schema\n--schema    Fetch Schema (also go to localhost:8080/v1/schema)\n--load      Load Data from ./init.json\n--query {a} Query DB for {a}\n")
		return
	}
	args := os.Args[1:]
	L:
		for i := 0; i < len(args); i++ {
			switch atmpt := args[i]; atmpt {
			case "--init":
				InitSchema(client)
			case "--delete":
				DeleteClass(client, "Recommendations")
			case "--schema":
				GetSchema(client)
			case "--load":
				LoadDataFromJson(client)
			case "--query":
				QueryData(client, args[i+1])
				break L
			default:
				fmt.Println(fmt.Sprintf("\narg %s not recognized, try one of the following: \n\n--init      Initialize Schema\n--delete    Delete Schema\n--schema    Fetch Schema (also go to localhost:8080/v1/schema)\n--load      Load Data from ./init.json\n--query {a} Query DB for {a}\n", atmpt))
			}
		}
}

func InitSchema(client *weaviate.Client) {
	check, err := client.Schema().ClassExistenceChecker().WithClassName("Recommendations").Do(context.Background())

	if err != nil {
		panic(err)
	}

	if !check {
		classObj := &models.Class{
			Class:       "Recommendations",
			Description: "Advice for how to represent yourself to recruiters and hiring managers",
			Vectorizer:  "text2vec-transformers", 
		}

		if err := client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()); 
		err != nil {
			panic(err)
		}


		positionProp := &models.Property{
			DataType: []string{"string"},
			Name: "DesiredPosition",
			Description: "The position the candidate is hoping to achieve",
		}

		sectionProp := &models.Property{
			DataType: []string{"string"},
			Name: "ProfileSection",
			Description: "The section of the candidate profile to recommend improvements on",
		}

		inputProp := &models.Property{
			DataType: []string{"string"},
			Name: "Input",
			Description: "The current input we want to give recommendations for how to improve",
		}

		suggestionProp := &models.Property{
			DataType: []string{"string"},
			Name: "Suggestion",
			Description: "Recommendation for how we can improve the input.",
		}

		proposedChangeProp := &models.Property{
			DataType: []string{"string"},
			Name: "ProposedChange",
			Description: "AI generated replacement.",
		}

		client.Schema().PropertyCreator().WithClassName("Recommendations").WithProperty(positionProp).Do(context.Background())
		client.Schema().PropertyCreator().WithClassName("Recommendations").WithProperty(sectionProp).Do(context.Background())
		client.Schema().PropertyCreator().WithClassName("Recommendations").WithProperty(inputProp).Do(context.Background())
		client.Schema().PropertyCreator().WithClassName("Recommendations").WithProperty(suggestionProp).Do(context.Background())
		client.Schema().PropertyCreator().WithClassName("Recommendations").WithProperty(proposedChangeProp).Do(context.Background())
		fmt.Println("Successfully initialized Recommendations Class")
	}
}


func GetSchema(client *weaviate.Client) {
    schema, err := client.Schema().Getter().Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v", schema)
}

func DeleteClass(client *weaviate.Client, className string) {

	check, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(context.Background())

	if err != nil{
		panic(err)
	}

	if check {
		err := client.Schema().ClassDeleter().WithClassName(className).Do(context.Background()); 
		if err!=nil{
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Successfully deleted class %s", className))
	}
}

func LoadDataFromJson(client *weaviate.Client) {
	data, err := os.Open("init.json")

	if err!=nil {
		panic(err)
	}


	var items []map[string]string

	if err:= json.NewDecoder(data).Decode(&items); err!=nil{
		panic(err)
	}

	// convert items into a slice of models.Object
	objects := make([]*models.Object, len(items))
	for i := range items {
	objects[i] = &models.Object{
		Class: "Recommendations",
		Properties: map[string]any{
			"desiredPosition": items[i]["desiredPosition"],
			"profileSection": items[i]["profileSection"],
			"input": items[i]["input"],
			"suggestion": items[i]["suggestion"],
			"proposedChange": items[i]["proposedChange"],
		},
	}
	}

	// batch write items
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
	panic(err)
	}
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			fmt.Println("batch load failed: ", res.Result.Errors.Error)
		}
	}
	fmt.Println("Successfully loaded data from JSON")
}


func QueryData(client *weaviate.Client, queryString string) {
	fields := []graphql.Field{
		{Name: "suggestion"},
		{Name: "proposedChange"},
	  }
	
	  nearText := client.GraphQL().NearTextArgBuilder().WithConcepts([]string{queryString})
	
	  result, err := client.GraphQL().Get().WithClassName("Recommendations").WithFields(fields...).WithNearText(nearText).WithLimit(1).Do(context.Background())

	  if err != nil {
		panic(err)
	  }
	

	  fmt.Println(result)
}