package robby

import (
	"math"
	"sort"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

type VectorRecord struct {
	Id               string    `json:"id"`
	Prompt           string    `json:"prompt"`
	Embedding        []float64 `json:"embedding"`
	CosineSimilarity float64
}

type MemoryVectorStore struct {
	Records map[string]VectorRecord
}

func (mvs *MemoryVectorStore) GetAll() ([]VectorRecord, error) {
	var records []VectorRecord
	for _, record := range mvs.Records {
		records = append(records, record)
	}
	return records, nil
}

// Save saves a vector record to the MemoryVectorStore.
// If the record does not have an ID, it generates a new UUID for it.
// It returns the saved vector record and an error if any occurred during the save operation.
// If the record already exists, it will be overwritten.
func (mvs *MemoryVectorStore) Save(vectorRecord VectorRecord) (VectorRecord, error) {
	if vectorRecord.Id == "" {
		vectorRecord.Id = uuid.New().String()
	}
	mvs.Records[vectorRecord.Id] = vectorRecord
	return vectorRecord, nil
}

// SearchSimilarities searches for vector records in the MemoryVectorStore that have a cosine distance similarity greater than or equal to the given limit.
//
// Parameters:
//   - embeddingFromQuestion: the vector record to compare similarities with.
//   - limit: the minimum cosine distance similarity threshold.
//
// Returns:
//   - []llm.VectorRecord: a slice of vector records that have a cosine distance similarity greater than or equal to the limit.
//   - error: an error if any occurred during the search.
func (mvs *MemoryVectorStore) SearchSimilarities(embeddingFromQuestion VectorRecord, limit float64) ([]VectorRecord, error) {

	var records []VectorRecord

	for _, v := range mvs.Records {
		distance := cosineSimilarity(embeddingFromQuestion.Embedding, v.Embedding)
		if distance >= limit {
			v.CosineSimilarity = distance
			records = append(records, v)
		}
	}
	return records, nil
}

// SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question.
// It returns a slice of vector records and an error if any.
// The limit parameter specifies the minimum similarity score for a record to be considered similar.
// The max parameter specifies the maximum number of vector records to return.
func (mvs *MemoryVectorStore) SearchTopNSimilarities(embeddingFromQuestion VectorRecord, limit float64, max int) ([]VectorRecord, error) {
	records, err := mvs.SearchSimilarities(embeddingFromQuestion, limit)
	if err != nil {
		return nil, err
	}
	return getTopNVectorRecords(records, max), nil
}

// getTopNVectorRecords returns the top N vector records based on their cosine similarity.
func getTopNVectorRecords(records []VectorRecord, max int) []VectorRecord {
	// Sort the records slice in descending order based on CosineDistance
	sort.Slice(records, func(i, j int) bool {
		return records[i].CosineSimilarity > records[j].CosineSimilarity
	})

	// Return the first max records or all if less than three
	if len(records) < max {
		return records
	}
	return records[:max]
}

// --- Cosine similarity ---

// dotProduct calculates the dot product of two vectors
// It assumes that both vectors are of the same length.
func dotProduct(v1 []float64, v2 []float64) float64 {
	// Calculate the dot product of two vectors
	sum := 0.0
	for i := range v1 {
		sum += v1[i] * v2[i]
	}
	return sum
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(v1, v2 []float64) float64 {
	// Calculate the cosine distance between two vectors
	product := dotProduct(v1, v2)

	norm1 := math.Sqrt(dotProduct(v1, v1))
	norm2 := math.Sqrt(dotProduct(v2, v2))
	if norm1 <= 0.0 || norm2 <= 0.0 {
		// Handle potential division by zero
		return 0.0
	}
	return product / (norm1 * norm2)
}

// RAGMemorySearchSimilaritiesWithText searches for similar records in the RAG memory using the provided text.
// It creates an embedding from the text and searches for records with cosine similarity above the specified limit.
// It returns a slice of strings containing the prompts of the similar records and an error if any occurred.
// If no similar records are found, it returns an empty slice.
// It requires the DMR client to be initialized and the embedding parameters to be set in the Agent.
// The limit parameter specifies the minimum cosine similarity score for a record to be considered similar.
// It returns an error if the embedding creation fails or if the search operation fails.
func (agent *Agent) RAGMemorySearchSimilaritiesWithText(text string, limit float64) ([]string, error) {
	// Create the embedding from the question
	agent.EmbeddingParams.Input = openai.EmbeddingNewParamsInputUnion{
		OfString: openai.String(text),
	}
	embeddingResponse, err := agent.dmrClient.Embeddings.New(agent.ctx, agent.EmbeddingParams)
	if err != nil {
		return nil, err
	}
	// -------------------------------------------------
	// Create a vector record from the user embedding
	// -------------------------------------------------
	embeddingFromText := VectorRecord{
		Embedding: embeddingResponse.Data[0].Embedding,
	}

	similarities, _ := agent.Store.SearchSimilarities(embeddingFromText, limit)
	var results []string
	for _, similarity := range similarities {
		results = append(results, similarity.Prompt)
	}
	return results, nil

}

// RAGMemorySearchSimilaritiesWith searches for similar records in the RAG memory using the provided embedding.
// It creates an embedding from the input and searches for records with cosine similarity above the specified limit.
// It returns a slice of strings containing the prompts of the similar records and an error if any occurred.
// If no similar records are found, it returns an empty slice.
// It requires the DMR client to be initialized and the embedding parameters to be set in the Agent.
// The limit parameter specifies the minimum cosine similarity score for a record to be considered similar.
// It returns an error if the embedding creation fails or if the search operation fails.
func (agent *Agent) RAGMemorySearchSimilaritiesWith(embedding openai.EmbeddingNewParamsInputUnion, limit float64) ([]string, error) {
		// Create the embedding from the question
	agent.EmbeddingParams.Input = embedding
	embeddingResponse, err := agent.dmrClient.Embeddings.New(agent.ctx, agent.EmbeddingParams)
	if err != nil {
		return nil, err
	}
	// -------------------------------------------------
	// Create a vector record from the user embedding
	// -------------------------------------------------
	embeddingFromText := VectorRecord{
		Embedding: embeddingResponse.Data[0].Embedding,
	}

	similarities, _ := agent.Store.SearchSimilarities(embeddingFromText, limit)
	var results []string
	for _, similarity := range similarities {
		results = append(results, similarity.Prompt)
	}
	return results, nil
}


// TODO: add helpers:
// - to create embeddings
// - then create a method: RAGMemorySearchSimilaritiesWith (embedding)