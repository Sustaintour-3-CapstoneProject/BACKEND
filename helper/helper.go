package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Struct untuk API Response
type Response struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

// APIResponse untuk membungkus response API
func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonresponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonresponse
}

// FormatValidationError memformat error validasi menjadi array string
func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}

// ValidateInput memvalidasi struct input menggunakan library validator
func ValidateInput(input interface{}) error {
	validate := validator.New()
	return validate.Struct(input)
}

// HashPassword mengenkripsi password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash mencocokkan password dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Struct untuk JWT Claims
type JWTClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat token JWT
func GenerateJWT(userID uint, username, role string) (string, error) {
	claims := &JWTClaims{
		Username: username,
		UserID:   userID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string  `json:"finishReason"`
		AvgLogprobs  float64 `json:"avgLogprobs"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	ModelVersion string `json:"modelVersion"`
}

func CallGeminiAPI(message string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	baseURL := os.Getenv("GEMINI_BASE_URL") + "?key=" + apiKey

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": message},
				},
			},
		},
	}

	if strings.Contains(message, "halo") {
		return "Selamat Datang di Tripwise, folks!", nil
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	// Handle the response
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response from API:")
		fmt.Println(string(body))
	} else {
		fmt.Println("Response from API:")
		fmt.Println(string(body))
	}

	var response GeminiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	latDiff := (lat2 - lat1) * (math.Pi / 180)
	lonDiff := (lon2 - lon1) * (math.Pi / 180)

	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)

	a := math.Sin(latDiff/2)*math.Sin(latDiff/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(lonDiff/2)*math.Sin(lonDiff/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
