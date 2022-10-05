package models

// TemplateData - holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // use interface as a return type in case you dont know what type of data you'll receive
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}
