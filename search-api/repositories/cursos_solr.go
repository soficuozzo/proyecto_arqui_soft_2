package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stevenferrer/solr-go"
	cursosDomain "proyecto_arqui_soft_2/search-api/domain"
)

type SolrConfig struct {
	Host       string // Solr host
	Port       string // Solr port
	Collection string // Solr collection name
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
	BaseURL    string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) Solr {
	// Construct the BaseURL using the provided host and port
    baseURL := fmt.Sprintf("http://%s:%s/solr/%s", config.Host, config.Port, config.Collection)
    client := solr.NewJSONClient(baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
		BaseURL:    baseURL,
	}
}

// Index adds a new cursos document to the Solr collection
func (searchEngine Solr) Index(ctx context.Context, curso cursosDomain.CursoData) (string, error) {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"curso_id": 	  curso.CursoID,
		"nombre": curso.Nombre,
		"descripcion": curso.Descripcion,
		"categoria": curso.Categoria,
		"capacidad": curso.Capacidad,
		"imagen": curso.Imagen,
		"valoracion": curso.Valoracion,
		"requisito": curso.Requisito,
		"profesor": curso.Profesor,
		"duracion": curso.Duracion,
	}

	// Prepare the index request
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Index the document in Solr
	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling curso document: %w", err)
	}
		// Log the URL being used for indexing
		fmt.Printf("Indexing URL: %s/update\n", searchEngine.BaseURL)

	// Index the document in Solr
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error indexing curso: %w", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("failed to index curso: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error committing changes to Solr: %w", err)
	}
	fmt.Println("Indexado")
	return curso.CursoID, nil
}

// Update modifies an existing curso document in the Solr collection
func (searchEngine Solr) Update(ctx context.Context, curso cursosDomain.CursoData) error {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"curso_id": 	  curso.CursoID,
		"nombre": curso.Nombre,
		"descripcion": curso.Descripcion,
		"categoria": curso.Categoria,
		"capacidad": curso.Capacidad,
		"imagen": curso.Imagen,
		"valoracion": curso.Valoracion,
		"requisito": curso.Requisito,
		"profesor": curso.Profesor,
		"duracion": curso.Duracion,
	}

	// Prepare the update request
	updateRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Update the document in Solr
	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error marshaling curso document: %w", err)
	}

	// Execute the update request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating curso: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to update curso: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Delete(ctx context.Context, id string) error {
	// Prepare the delete request
	docToDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	// Update the document in Solr
	body, err := json.Marshal(docToDelete)
	if err != nil {
		return fmt.Errorf("error marshaling curso document: %w", err)
	}

	// Execute the delete request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error deleting curso: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to index curso: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]cursosDomain.CursoData, error) {
	// Prepare the Solr query with limit and offset
	solrQuery := fmt.Sprintf("q=(nombre:%s)&rows=%d&start=%d", query, limit, offset)
	fmt.Printf("Searching URL: %s/select?%s\n", searchEngine.BaseURL, solrQuery)

	// Execute the search request
	resp, err := searchEngine.Client.Query(ctx, searchEngine.Collection, solr.NewQuery(solrQuery))
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to execute search query: %v", resp.Error)
	}

	// Parse the response and extract curso documents
	var cursosList []cursosDomain.CursoData
	for _, doc := range resp.Response.Documents {
	

		// Safely extract curso fields with type assertions
		cursos := cursosDomain.CursoData{
			CursoID:        getStringField(doc, "curso_id"),
			Nombre:      getStringField(doc, "nombre"),
			Descripcion:   getStringField(doc, "descripcion"),
			Categoria:      getStringField(doc, "categoria"),
			Capacidad:     getIntField(doc, "capacidad"),
			Imagen:      getStringField(doc, "imagen"),
			Valoracion:  getIntField(doc, "valoracion"),
			Requisito:   getStringField(doc, "requisito"),
			Profesor:   getStringField(doc, "profesor"),
			Duracion:  	getIntField(doc, "duracion"),
		}
		cursosList = append(cursosList, cursos)
	}

	return cursosList, nil
}



// Helper function to safely get string fields from the document
func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			return strVal
		}
	}
	return ""
}


// Helper function to safely get int64 fields from the document
func getIntField(doc map[string]interface{}, field string) int64 {
	if val, ok := doc[field].(int64); ok {
		return val
	}
	if val, ok := doc[field].(float64); ok {
		return int64(val)
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if intVal, ok := val[0].(int64); ok {
			return intVal
		}
		if floatVal, ok := val[0].(float64); ok {
			return int64(floatVal)
		}
	}
	return 0
}
