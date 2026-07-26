package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/nutrixpos/pos/common"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	aimodels "github.com/nutrixpos/pos/modules/ai/models"
	coremodels "github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AIService struct {
	Logger logger.ILogger
	Config config.Config
}

func normalizeSlovenian(s string) string {
	replacer := strings.NewReplacer(
		"č", "c", "Č", "C",
		"š", "s", "Š", "S",
		"ž", "z", "Ž", "Z",
	)
	return replacer.Replace(s)
}

func soundex(s string) string {
	s = normalizeSlovenian(strings.ToLower(s))
	if len(s) == 0 {
		return ""
	}

	var result []byte
	result = append(result, s[0])

	codes := map[byte]byte{
		'b': '1', 'f': '1', 'p': '1', 'v': '1',
		'c': '2', 'g': '2', 'j': '2', 'k': '2', 'q': '2', 's': '2', 'x': '2', 'z': '2',
		'd': '3', 't': '3',
		'l': '4',
		'm': '5', 'n': '5',
		'r': '6',
	}

	var prev byte
	for i := 1; i < len(s) && len(result) < 4; i++ {
		code, ok := codes[s[i]]
		if ok && code != prev {
			result = append(result, code)
			prev = code
		} else if !ok {
			prev = 0
		}
	}

	for len(result) < 4 {
		result = append(result, '0')
	}

	return string(result[:4])
}

func levenshteinDistance(s1, s2 string) int {
	s1 = normalizeSlovenian(strings.ToLower(s1))
	s2 = normalizeSlovenian(strings.ToLower(s2))

	if s1 == s2 {
		return 0
	}
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	prev := make([]int, len(s2)+1)
	curr := make([]int, len(s2)+1)

	for j := 0; j <= len(s2); j++ {
		prev[j] = j
	}

	for i := 1; i <= len(s1); i++ {
		curr[0] = i
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			curr[j] = min3(
				prev[j]+1,
				curr[j-1]+1,
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[len(s2)]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func FuzzyMatch(query, target string) float64 {
	q := normalizeSlovenian(strings.ToLower(strings.TrimSpace(query)))
	t := normalizeSlovenian(strings.ToLower(strings.TrimSpace(target)))

	if len(q) == 0 {
		return 0
	}
	if len(t) == 0 {
		return 0
	}

	if q == t {
		return 1.0
	}

	if strings.HasPrefix(t, q) {
		return 0.9
	}

	if strings.Contains(t, q) {
		return 0.8
	}

	if strings.HasPrefix(q, t) {
		return 0.75
	}

	qWords := strings.Fields(q)
	tWords := strings.Fields(t)
	if len(qWords) > 1 || len(tWords) > 1 {
		matchCount := 0
		for _, qw := range qWords {
			for _, tw := range tWords {
				if strings.Contains(tw, qw) || strings.Contains(qw, tw) {
					matchCount++
					break
				}
			}
		}
		if matchCount > 0 {
			ratio := float64(matchCount) / float64(len(qWords))
			return 0.5 + 0.2*ratio
		}
	}

	qSoundex := soundex(q)
	tSoundex := soundex(t)
	if qSoundex == tSoundex && len(qSoundex) > 0 {
		return 0.5
	}

	distance := levenshteinDistance(q, t)
	maxLen := float64(len(q))
	if float64(len(t)) > maxLen {
		maxLen = float64(len(t))
	}
	if maxLen == 0 {
		return 0
	}

	similarity := 1.0 - float64(distance)/maxLen

	if similarity > 0.7 {
		return 0.4 + 0.3*similarity
	}

	if similarity > 0.5 {
		return 0.2 + 0.2*similarity
	}

	return similarity * 0.2
}

func (s *AIService) AISearch(query string, branchId string, language string, limit int) (*aimodels.AISearchResponse, error) {
	if limit <= 0 {
		limit = 20
	}

	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		return nil, fmt.Errorf("AISearch: %w", err)
	}

	ctx := context.Background()
	dbName := s.Config.Databases[0].Database

	productsColl := client.Database(dbName).Collection("recipes")
	cursor, err := productsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("AISearch: %w", err)
	}
	defer cursor.Close(ctx)

	var products []coremodels.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("AISearch: %w", err)
	}

	categoriesColl := client.Database(dbName).Collection("categories")
	catCursor, err := categoriesColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("AISearch: %w", err)
	}
	defer catCursor.Close(ctx)

	var categories []coremodels.Category
	if err := catCursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("AISearch: %w", err)
	}

	productCategoryMap := make(map[string]string)
	for _, cat := range categories {
		for _, p := range cat.Products {
			productCategoryMap[p.Id] = cat.Name
		}
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))
	queryNormalized := normalizeSlovenian(queryLower)
	queryWords := strings.Fields(queryLower)

	type scoredResult struct {
		result aimodels.SearchResult
		score  float64
	}

	var scored []scoredResult

	for _, product := range products {
		productNameLower := strings.ToLower(product.Name)
		productNameNormalized := normalizeSlovenian(productNameLower)

		score := 0.0

		if productNameLower == queryLower || productNameNormalized == queryNormalized {
			score = 100.0
		} else if strings.HasPrefix(productNameNormalized, queryNormalized) {
			score = 80.0
		} else if strings.Contains(productNameNormalized, queryNormalized) {
			score = 60.0
		} else if strings.HasPrefix(queryNormalized, productNameNormalized) {
			score = 55.0
		} else {
			fuzzyScore := FuzzyMatch(query, product.Name)
			if fuzzyScore > 0.5 {
				score = 40.0 + 19.0*fuzzyScore
			} else if fuzzyScore > 0.3 {
				score = 20.0 + 20.0*fuzzyScore
			}
		}

		if catName, ok := productCategoryMap[product.Id]; ok {
			catLower := strings.ToLower(catName)
			if strings.Contains(catLower, queryLower) || strings.Contains(queryLower, catLower) {
				if score < 30.0 {
					score = 30.0
				}
			}
		}

		for _, material := range product.Materials {
			matLower := strings.ToLower(material.Name)
			matNormalized := normalizeSlovenian(matLower)
			if strings.Contains(matNormalized, queryNormalized) || strings.Contains(queryNormalized, matNormalized) {
				if score < 20.0 {
					score = 20.0
				}
			}
		}

		if len(queryWords) > 1 {
			multiWordScore := 0.0
			allMatched := true
			for _, qw := range queryWords {
				qwNorm := normalizeSlovenian(qw)
				wordFound := false
				if strings.Contains(productNameNormalized, qwNorm) {
					wordFound = true
					multiWordScore += 50.0
				} else {
					for _, mat := range product.Materials {
						matNorm := normalizeSlovenian(strings.ToLower(mat.Name))
						if strings.Contains(matNorm, qwNorm) {
							wordFound = true
							multiWordScore += 20.0
							break
						}
					}
				}
				if !wordFound {
					allMatched = false
				}
			}
			if allMatched && multiWordScore > score {
				score = multiWordScore
			}
		}

		if score > 0 {
			desc := fmt.Sprintf("%s - %.2f EUR", product.Name, product.Price)

			scored = append(scored, scoredResult{
				result: aimodels.SearchResult{
					ProductId:   product.Id,
					Name:        product.Name,
					Description: desc,
					Price:       product.Price,
					Score:       score,
					Category:    productCategoryMap[product.Id],
				},
				score: score,
			})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if limit > len(scored) {
		limit = len(scored)
	}

	results := make([]aimodels.SearchResult, limit)
	for i := 0; i < limit; i++ {
		results[i] = scored[i].result
	}

	suggestions := s.getPopularSuggestions(branchId, 5)

	return &aimodels.AISearchResponse{
		Results:     results,
		Suggestions: suggestions,
	}, nil
}

func (s *AIService) ProcessVoiceOrder(audioBase64 string, language string, branchId string) (*aimodels.VoiceOrderResponse, error) {
	if audioBase64 == "" {
		return nil, fmt.Errorf("ProcessVoiceOrder: audio_base64 is required")
	}

	_, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return nil, fmt.Errorf("ProcessVoiceOrder: invalid base64 audio: %w", err)
	}

	transcript := s.extractKeywordsFromAudio(audioBase64, language)

	if transcript == "" {
		transcript = "order"
	}

	items, confidence := s.matchTranscriptToItems(transcript, branchId)

	suggestions, _ := s.GetSmartSuggestions(branchId, "", 5)
	suggestionTexts := make([]string, 0, len(suggestions))
	for _, sug := range suggestions {
		suggestionTexts = append(suggestionTexts, sug.ProductName)
	}

	return &aimodels.VoiceOrderResponse{
		Transcript:  transcript,
		Items:       items,
		Confidence:  confidence,
		Suggestions: suggestionTexts,
	}, nil
}

func (s *AIService) extractKeywordsFromAudio(audioBase64 string, language string) string {
	if len(audioBase64) < 100 {
		return "order"
	}

	sampleBytes, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return "order"
	}

	var keywords []byte
	step := 1
	if len(sampleBytes) > 200 {
		step = len(sampleBytes) / 100
	}
	for i := 0; i < len(sampleBytes) && len(keywords) < 100; i += step {
		b := sampleBytes[i]
		if b >= 32 && b <= 126 {
			keywords = append(keywords, b)
		} else if b == 32 || b == 10 || b == 13 {
			keywords = append(keywords, ' ')
		}
	}

	result := strings.TrimSpace(string(keywords))
	if result == "" {
		return "order"
	}

	return result
}

func (s *AIService) matchTranscriptToItems(transcript string, branchId string) ([]aimodels.AIOrderItem, float64) {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("matchTranscriptToItems: %s", err.Error()))
		return nil, 0
	}

	ctx := context.Background()
	dbName := s.Config.Databases[0].Database

	productsColl := client.Database(dbName).Collection("recipes")
	cursor, err := productsColl.Find(ctx, bson.M{})
	if err != nil {
		s.Logger.Error(fmt.Sprintf("matchTranscriptToItems: %s", err.Error()))
		return nil, 0
	}
	defer cursor.Close(ctx)

	var products []coremodels.Product
	if err := cursor.All(ctx, &products); err != nil {
		s.Logger.Error(fmt.Sprintf("matchTranscriptToItems: %s", err.Error()))
		return nil, 0
	}

	transcriptLower := strings.ToLower(transcript)
	words := strings.Fields(transcriptLower)

	quantityWords := map[string]bool{
		"one": true, "ena": true, "jedan": true,
		"two": true, "dve": true, "dva": true,
		"three": true, "tri": true, "štiri": true, "četiri": true,
		"five": true, "pet": true,
		"en": true, "sl": true, "hr": true,
	}

	var items []aimodels.AIOrderItem
	usedProducts := make(map[string]bool)

	for _, word := range words {
		if quantityWords[word] {
			continue
		}

		bestScore := 0.0
		var bestProduct *coremodels.Product

		for i := range products {
			if usedProducts[products[i].Id] {
				continue
			}

			score := FuzzyMatch(word, products[i].Name)
			if score > bestScore {
				bestScore = score
				bestProduct = &products[i]
			}

			for _, mat := range products[i].Materials {
				matScore := FuzzyMatch(word, mat.Name)
				if matScore > bestScore {
					bestScore = matScore
					bestProduct = &products[i]
				}
			}
		}

		if bestProduct != nil && bestScore > 0.4 {
			items = append(items, aimodels.AIOrderItem{
				ProductName: bestProduct.Name,
				Quantity:    1.0,
				Comment:     "",
				Confidence:  bestScore,
				ProductId:   bestProduct.Id,
			})
			usedProducts[bestProduct.Id] = true
		}
	}

	totalConfidence := 0.0
	if len(items) > 0 {
		for _, item := range items {
			totalConfidence += item.Confidence
		}
		totalConfidence /= float64(len(items))
	}

	return items, totalConfidence
}

func (s *AIService) GetSmartSuggestions(branchId string, orderId string, limit int) ([]aimodels.SmartSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}

	suggestions := s.getTimeBasedSuggestions(branchId, limit)

	timeSuggestions := s.getTimeBasedSuggestions(branchId, limit)
	for _, ts := range timeSuggestions {
		found := false
		for i, existing := range suggestions {
			if existing.ProductId == ts.ProductId {
				suggestions[i].Score += ts.Score
				found = true
				break
			}
		}
		if !found {
			suggestions = append(suggestions, ts)
		}
	}

	popularSuggestions := s.getPopularSuggestions(branchId, limit)
	for _, ps := range popularSuggestions {
		found := false
		for i, existing := range suggestions {
			if existing.ProductId == ps.ProductId {
				suggestions[i].Score += ps.Score
				found = true
				break
			}
		}
		if !found {
			suggestions = append(suggestions, ps)
		}
	}

	seasonalSuggestions := s.getSeasonalSuggestions(limit)
	for _, ss := range seasonalSuggestions {
		found := false
		for i, existing := range suggestions {
			if existing.ProductId == ss.ProductId {
				suggestions[i].Score += ss.Score
				found = true
				break
			}
		}
		if !found {
			suggestions = append(suggestions, ss)
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	if limit > len(suggestions) {
		limit = len(suggestions)
	}

	return suggestions[:limit], nil
}

func (s *AIService) getTimeBasedSuggestions(branchId string, limit int) []aimodels.SmartSuggestion {
	hour := time.Now().Hour()

	categoryKeywords := []string{}

	switch {
	case hour >= 6 && hour < 10:
		categoryKeywords = []string{"zajtrk", "breakfast", "kava", "coffee", "čaj", "tea", "croissant", "sendvič", "sandwich"}
	case hour >= 11 && hour < 14:
		categoryKeywords = []string{"kosilo", "lunch", "solata", "salad", "juha", "soup", "glavna", "main"}
	case hour >= 17 && hour < 21:
		categoryKeywords = []string{"večerja", "dinner", "pizza", "testenine", "pasta", "meso", "meat", "rib", "fish"}
	default:
		return nil
	}

	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("getTimeBasedSuggestions: %s", err.Error()))
		return nil
	}

	ctx := context.Background()
	dbName := s.Config.Databases[0].Database

	productsColl := client.Database(dbName).Collection("recipes")
	categoriesColl := client.Database(dbName).Collection("categories")

	catCursor, err := categoriesColl.Find(ctx, bson.M{})
	if err != nil {
		s.Logger.Error(fmt.Sprintf("getTimeBasedSuggestions: %s", err.Error()))
		return nil
	}
	defer catCursor.Close(ctx)

	var categories []coremodels.Category
	if err := catCursor.All(ctx, &categories); err != nil {
		s.Logger.Error(fmt.Sprintf("getTimeBasedSuggestions: %s", err.Error()))
		return nil
	}

	var matchingProductIds []string
	for _, cat := range categories {
		catLower := strings.ToLower(cat.Name)
		catNorm := normalizeSlovenian(catLower)
		for _, kw := range categoryKeywords {
			kwNorm := normalizeSlovenian(strings.ToLower(kw))
			if strings.Contains(catNorm, kwNorm) {
				for _, p := range cat.Products {
					matchingProductIds = append(matchingProductIds, p.Id)
				}
			}
		}
	}

	if len(matchingProductIds) == 0 {
		return nil
	}

	var suggestions []aimodels.SmartSuggestion
	for _, pid := range matchingProductIds {
		var product coremodels.Product
		err := productsColl.FindOne(ctx, bson.M{"id": pid}).Decode(&product)
		if err != nil {
			continue
		}

		suggestions = append(suggestions, aimodels.SmartSuggestion{
			ProductId:   product.Id,
			ProductName: product.Name,
			Reason:      "time_based",
			Score:       30.0,
		})
	}

	if limit > len(suggestions) {
		limit = len(suggestions)
	}
	if limit > 0 {
		return suggestions[:limit]
	}
	return suggestions
}

func (s *AIService) getPopularSuggestions(branchId string, limit int) []aimodels.SmartSuggestion {
	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("getPopularSuggestions: %s", err.Error()))
		return nil
	}

	ctx := context.Background()
	dbName := s.Config.Databases[0].Database

	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	pipeline := []bson.M{
		{"$match": bson.M{
			"submitted_at": bson.M{"$gte": oneWeekAgo},
		}},
		{"$unwind": "$items"},
		{"$group": bson.M{
			"_id":        "$items.product.id",
			"orderCount": bson.M{"$sum": 1},
			"totalQty":   bson.M{"$sum": "$items.quantity"},
		}},
		{"$sort": bson.M{"totalQty": -1}},
		{"$limit": int64(limit * 2)},
	}

	ordersColl := client.Database(dbName).Collection("orders")
	cursor, err := ordersColl.Aggregate(ctx, pipeline)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("getPopularSuggestions: %s", err.Error()))
		return nil
	}
	defer cursor.Close(ctx)

	type popularityResult struct {
		Id         string  `bson:"_id"`
		OrderCount int     `bson:"orderCount"`
		TotalQty   float64 `bson:"totalQty"`
	}

	var results []popularityResult
	if err := cursor.All(ctx, &results); err != nil {
		s.Logger.Error(fmt.Sprintf("getPopularSuggestions: %s", err.Error()))
		return nil
	}

	productsColl := client.Database(dbName).Collection("recipes")

	var suggestions []aimodels.SmartSuggestion
	for _, r := range results {
		if r.Id == "" {
			continue
		}

		var product coremodels.Product
		err := productsColl.FindOne(ctx, bson.M{"id": r.Id}).Decode(&product)
		if err != nil {
			continue
		}

		score := float64(r.OrderCount) * 5.0
		if score > 50.0 {
			score = 50.0
		}

		suggestions = append(suggestions, aimodels.SmartSuggestion{
			ProductId:   product.Id,
			ProductName: product.Name,
			Reason:      "popular",
			Score:       score,
		})
	}

	return suggestions
}

func (s *AIService) getSeasonalSuggestions(limit int) []aimodels.SmartSuggestion {
	month := time.Now().Month()

	seasonalKeywords := []string{}

	switch {
	case month >= 3 && month <= 5:
		seasonalKeywords = []string{"pomlad", "spring", "solata", "salad", "sveže", "fresh"}
	case month >= 6 && month <= 8:
		seasonalKeywords = []string{"poletje", "summer", "led", "ice", "sadje", "fruit", "grilled"}
	case month >= 9 && month <= 11:
		seasonalKeywords = []string{"jesen", "autumn", "gobe", "mushroom", "buče", "pumpkin"}
	default:
		seasonalKeywords = []string{"zima", "winter", "toplo", "warm", "juha", "soup"}
	}

	client, err := common.GetDatabaseClient(s.Logger, &s.Config)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("getSeasonalSuggestions: %s", err.Error()))
		return nil
	}

	ctx := context.Background()
	dbName := s.Config.Databases[0].Database

	productsColl := client.Database(dbName).Collection("recipes")
	cursor, err := productsColl.Find(ctx, bson.M{})
	if err != nil {
		s.Logger.Error(fmt.Sprintf("getSeasonalSuggestions: %s", err.Error()))
		return nil
	}
	defer cursor.Close(ctx)

	var products []coremodels.Product
	if err := cursor.All(ctx, &products); err != nil {
		s.Logger.Error(fmt.Sprintf("getSeasonalSuggestions: %s", err.Error()))
		return nil
	}

	var suggestions []aimodels.SmartSuggestion
	for _, product := range products {
		productLower := strings.ToLower(product.Name)
		productNorm := normalizeSlovenian(productLower)

		for _, kw := range seasonalKeywords {
			kwNorm := normalizeSlovenian(strings.ToLower(kw))
			if strings.Contains(productNorm, kwNorm) {
				suggestions = append(suggestions, aimodels.SmartSuggestion{
					ProductId:   product.Id,
					ProductName: product.Name,
					Reason:      "seasonal",
					Score:       25.0,
				})
				break
			}
		}
	}

	if limit > len(suggestions) {
		limit = len(suggestions)
	}
	if limit > 0 {
		return suggestions[:limit]
	}
	return suggestions
}

var _ = regexp.MustCompile
