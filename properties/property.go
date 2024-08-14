package properties

type Metadata = map[string]interface{}
type Properties = map[string]interface{}
type Collections = map[string]Properties
type Data struct {
	Properties  Properties  `json:"properties" dynamodbav:"properties"`
	Collections Collections `json:"collections" dynamodbav:"collections"`
}
