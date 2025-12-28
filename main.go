package main

// Import the required packages
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/server"
)

// Define Input and Output Structures for each flow

// BloodSugar Input Struct
type BloodSugarInput struct {
	Reading    float64 `json:"reading" jsonschema:"description=Blood sugar reading in mg/dL"`
	MealTiming string  `json:"meal_timing" jsonschema:"description=Timing: fasting, before_meal, after_meal"`
	MealType   string  `json:"meal_type" jsonschema:"description=Type of meal: breakfast, lunch, dinner, snack"`
}

// BloodSugar Output Struct
type BloodSugarOutput struct {
	Status         string `json:"status" jsonschema:"description=Status: normal, high, low, critical"`
	Interpretation string `json:"interpretation" jsonschema:"description=Detailed interpretation"`
	Recommendation string `json:"recommendation" jsonschema:"description=Immediate recommendations"`
}

// MealPlan Input Struct
type MealPlanInput struct {
	DietType     string  `json:"diet_type" jsonschema:"description=Diet preference: vegetarian, non_vegetarian, vegan"`
	Allergies    string  `json:"allergies" jsonschema:"description=Any food allergies or restrictions"`
	CalorieLimit float64 `json:"calorie_limit" jsonschema:"description=Daily calorie limit (optional)"`
}

// MealPlan Output Struct
type MealPlanOutput struct {
	Breakfast string `json:"breakfast" jsonschema:"description=Breakfast suggestions"`
	Lunch     string `json:"lunch" jsonschema:"description=Lunch suggestions"`
	Dinner    string `json:"dinner" jsonschema:"description=Dinner suggestions"`
	Snacks    string `json:"snacks" jsonschema:"description=Healthy snack options"`
}

// Symptom Input Struct
type SymptomInput struct {
	Symptoms    string `json:"symptoms" jsonschema:"description=Describe symptoms you're experiencing"`
	Duration    string `json:"duration" jsonschema:"description=How long symptoms have been present"`
	CurrentMeds string `json:"current_meds" jsonschema:"description=Current medications (optional)"`
}

// Symptom Output Struct
type SymptomOutput struct {
	Urgency    string `json:"urgency" jsonschema:"description=Urgency level: emergency, urgent, routine"`
	Assessment string `json:"assessment" jsonschema:"description=Symptom assessment"`
	NextSteps  string `json:"next_steps" jsonschema:"description=Recommended next steps"`
}

// Exercise Input Struct
type ExerciseInput struct {
	FitnessLevel  string  `json:"fitness_level" jsonschema:"description=Fitness level: beginner, intermediate, advanced"`
	TimeAvailable int     `json:"time_available" jsonschema:"description=Minutes available for exercise"`
	CurrentBG     float64 `json:"current_bg" jsonschema:"description=Current blood glucose level (optional)"`
	PreferredType string  `json:"preferred_type" jsonschema:"description=Exercise preference: cardio, strength, yoga, walking"`
}

// Exercise Output Struct
type ExerciseOutput struct {
	SafetyCheck    string `json:"safety_check" jsonschema:"description=Safety considerations based on BG"`
	Recommendation string `json:"recommendation" jsonschema:"description=Exercise recommendations"`
	Duration       string `json:"duration" jsonschema:"description=Recommended duration and intensity"`
	Precautions    string `json:"precautions" jsonschema:"description=Important precautions"`
}

// Medication Input Struct
type MedicationInput struct {
	MedicationName string `json:"medication_name" jsonschema:"description=Name of medication"`
	Purpose        string `json:"purpose" jsonschema:"description=Purpose of inquiry (dosage, timing, side_effects, interactions)"`
}

// Medication Output Struct
type MedicationOutput struct {
	Information string `json:"information" jsonschema:"description=Medication information"`
	Reminder    string `json:"reminder" jsonschema:"description=Important reminders"`
	Disclaimer  string `json:"disclaimer" jsonschema:"description=Medical disclaimer"`
}

// Declare main function
func main() {

	// Create a blank context
	ctx := context.Background()

	// Load the Google's AI API Key environment variable
	apiKey := os.Getenve("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI API KEY environment variable is missing!")
	}

	// Initialize Google's AI plugin with the Key
	plugin := &googlegenai.GoogleAI{
		APIKey: apiKey,
	}

	// Initialize Genkit
	g := genkit.Init(ctx,
		genkit.WithPlugins(plugin),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)

	// Welcome Message
	fmt.Println("=== DiabetesAI Advisor Initializing ===")
	response, err := genkit.Generate(ctx, g,
		ai.WithPrompt("Generate a warm welcome, encouraging welcome message for diabetes patients using this AI health advisor. Keep it under 50 words."),
	)
	if err != nil {
		log.Printf("Error generating welcome: %v", err)
	} else {
		fmt.Println("\n" + response.Text())
	}

	// Flow 1: Blood Sugar Interpreter
	bloodSugarFlow := genkit.DefineFlow(g, "bloodSugarInterpreter", func(ctx context.Context, input *BloodSugarInput) (*BloodSugarOutput, error) {
		prompt := fmt.Sprintf(`You are a diabetes care advisor. Analyze this blood sugar reading:
		
Reading: %.1f mg/dL
Timing: %s
Meal: %s

Provide:
1. Status (normal/high/low/critical)
2. Clear interpretation in simple terms
3. Immediate actionable recommendations

Guidelines:
- Fasting: 70-100 normal, 100-126 pre-diabetes, >126 diabetes concern
- Before meal: 70-130 normal
- 2 hours after meal: <180 normal
- <70 is low (hypoglycemia)
- >250 requires immediate attention

Be supportive and clear.`, input.Reading, input.MealTiming, input.MealType)

		result, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt))
		if err != nil {
			return nil, fmt.Errorf("failed to interpret blood sugar: %w", err)
		}

		// Determine status based on reading
		status := "normal"
		if input.Reading < 70 {
			status = "low"
		} else if input.Reading > 250 {
			status = "critical"
		} else if input.Reading > 180 {
			status = "high"
		}

		text := result.Text()
		parts := splitIntoSections(text, 3)

		return &BloodSugarOutput{
			Status:         status,
			Interpretation: parts[0],
			Recommendation: parts[1],
		}, nil
	})

	// Flow 2: Meal Planner
	mealPlanFlow := genkit.DefineFlow(g, "mealPlanner", func(ctx context.Context, input *MealPlanInput) (*MealPlanOutput, error) {
		calorieInfo := ""
		if input.CalorieLimit > 0 {
			calorieInfo = fmt.Sprintf("Target daily calories: %.0f", input.CalorieLimit)
		}

		prompt := fmt.Sprintf(`Create a diabetes-friendly meal plan:

Diet Type: %s
Allergies/Restrictions: %s
%s

For each meal, provide:
- Specific food items
- Approximate portion sizes
- Why it's good for blood sugar control

Focus on:
- Low glycemic index foods
- Balanced macros (protein, healthy fats, complex carbs)
- High fiber content
- Foods that prevent blood sugar spikes

Format:
BREAKFAST: [meal details]
LUNCH: [meal details]
DINNER: [meal details]
SNACKS: [snack options]`, input.DietType, input.Allergies, calorieInfo)

		result, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt))
		if err != nil {
			return nil, fmt.Errorf("failed to generate meal plan: %w", err)
		}

		text := result.Text()
		sections := parseMealSections(text)

		return &MealPlanOutput{
			Breakfast: sections["breakfast"],
			Lunch:     sections["lunch"],
			Dinner:    sections["dinner"],
			Snacks:    sections["snacks"],
		}, nil
	})

	// Flow 3: Symptom Checker
	symptomFlow := genkit.DefineFlow(g, "symptomChecker", func(ctx context.Context, input *SymptomInput) (*SymptomOutput, error) {
		prompt := fmt.Sprintf(`You are a diabetes health advisor. Assess these symptoms:

Symptoms: %s
Duration: %s
Current Medications: %s

Determine:
1. URGENCY LEVEL: 
   - EMERGENCY (call 911): Severe symptoms like chest pain, loss of consciousness, extreme confusion
   - URGENT (contact doctor today): Persistent high BG, signs of infection, concerning symptoms
   - ROUTINE (monitor and schedule appointment): Mild symptoms

2. ASSESSMENT: What these symptoms might indicate

3. NEXT STEPS: Specific actions to take

Be clear about when to seek immediate medical help. Always err on the side of caution.`, input.Symptoms, input.Duration, input.CurrentMeds)

		result, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt))
		if err != nil {
			return nil, fmt.Errorf("failed to check symptoms: %w", err)
		}

		text := result.Text()

		// Determine urgency from response
		urgency := "routine"
		if containsKeywords(text, []string{"emergency", "911", "immediate", "urgent care"}) {
			urgency = "emergency"
		} else if containsKeywords(text, []string{"urgent", "contact doctor", "today"}) {
			urgency = "urgent"
		}

		parts := splitIntoSections(text, 3)

		return &SymptomOutput{
			Urgency:    urgency,
			Assessment: parts[0],
			NextSteps:  parts[1],
		}, nil
	})

	// Flow 4: Exercise Advisor
	exerciseFlow := genkit.DefineFlow(g, "exerciseAdvisor", func(ctx context.Context, input *ExerciseInput) (*ExerciseOutput, error) {
		bgInfo := ""
		if input.CurrentBG > 0 {
			bgInfo = fmt.Sprintf("Current Blood Glucose: %.1f mg/dL", input.CurrentBG)
		}

		prompt := fmt.Sprintf(`Create a diabetes-safe exercise plan:

Fitness Level: %s
Time Available: %d minutes
%s
Preferred Exercise: %s

Provide:
1. SAFETY CHECK: Is it safe to exercise now based on BG? (BG 100-250 is generally safe, <100 eat snack first, >250 delay exercise)
2. EXERCISE PLAN: Specific exercises with sets/reps or duration
3. DURATION & INTENSITY: How to structure the workout
4. PRECAUTIONS: Important safety tips

Remember:
- Exercise lowers blood sugar
- Stay hydrated
- Have fast-acting carbs nearby
- Stop if feeling dizzy or unwell`, input.FitnessLevel, input.TimeAvailable, bgInfo, input.PreferredType)

		result, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt))
		if err != nil {
			return nil, fmt.Errorf("failed to generate exercise plan: %w", err)
		}

		text := result.Text()
		parts := splitIntoSections(text, 4)

		return &ExerciseOutput{
			SafetyCheck:    parts[0],
			Recommendation: parts[1],
			Duration:       parts[2],
			Precautions:    parts[3],
		}, nil
	})

	// Flow 5: Medication Info
	medicationFlow := genkit.DefineFlow(g, "medicationInfo", func(ctx context.Context, input *MedicationInput) (*MedicationOutput, error) {
		prompt := fmt.Sprintf(`Provide general information about diabetes medication:

Medication: %s
Question about: %s

Provide helpful general information, but:
1. DO NOT prescribe or change dosages
2. Emphasize consulting with healthcare provider
3. Mention common considerations
4. Include important safety information

Always include a clear disclaimer that this is educational information only.`, input.MedicationName, input.Purpose)

		result, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt))
		if err != nil {
			return nil, fmt.Errorf("failed to get medication info: %w", err)
		}

		disclaimer := "⚠️ IMPORTANT: This is educational information only. Always consult your healthcare provider before starting, stopping, or changing any medication. This AI advisor cannot replace professional medical advice."

		return &MedicationOutput{
			Information: result.Text(),
			Reminder:    "Set reminders on your phone for medication times. Never skip doses without consulting your doctor.",
			Disclaimer:  disclaimer,
		}, nil
	})

	// Flows' local tests
	fmt.Println("\n=== Testing Blood Sugar Interpreter ===")
	bsResult, err := bloodSugarFlow.Run(ctx, &BloodSugarInput{
		Reading:    145,
		MealTiming: "after_meal",
		MealType:   "lunch",
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Status: %s\n", bsResult.Status)
		fmt.Printf("Interpretation: %s\n", bsResult.Interpretation)
	}

	fmt.Println("\n=== Testing Meal Planner ===")
	mealResult, err := mealPlanFlow.Run(ctx, &MealPlanInput{
		DietType:     "vegetarian",
		Allergies:    "none",
		CalorieLimit: 1800,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Breakfast: %s\n", mealResult.Breakfast[:100]+"...")
	}

	fmt.Println("\n=== Testing Exercise Advisor ===")
	exerciseResult, err := exerciseFlow.Run(ctx, &ExerciseInput{
		FitnessLevel:  "beginner",
		TimeAvailable: 30,
		CurrentBG:     120,
		PreferredType: "walking",
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Safety Check: %s\n", exerciseResult.SafetyCheck)
	}

	// Set up HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bloodSugar", genkit.Handler(bloodSugarFlow))
	mux.HandleFunc("POST /mealPlan", genkit.Handler(mealPlanFlow))
	mux.HandleFunc("POST /symptoms", genkit.Handler(symptomFlow))
	mux.HandleFunc("POST /exercise", genkit.Handler(exerciseFlow))
	mux.HandleFunc("POST /medication", genkit.Handler(medicationFlow))

	// Print server info
	fmt.Println("\n=== DiabetesAI Advisor Server Starting ===")
	fmt.Println("Server: http://localhost:3400")
	fmt.Println("\nAvailable Endpoints:")
	fmt.Println("  POST /bloodSugar   - Interpret blood sugar readings")
	fmt.Println("  POST /mealPlan     - Get diabetes-friendly meal plans")
	fmt.Println("  POST /symptoms     - Check symptoms and get guidance")
	fmt.Println("  POST /exercise     - Get safe exercise recommendations")
	fmt.Println("  POST /medication   - Get medication information")
	fmt.Println("\nSample curl command:")
	fmt.Println(`  curl -X POST "http://localhost:3400/bloodSugar" \`)
	fmt.Println(`    -H "Content-Type: application/json" \`)
	fmt.Println(`    -d '{"data": {"reading": 145, "meal_timing": "after_meal", "meal_type": "lunch"}}'`)

	// Start the server
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}

// Helper function to split text into sections
func splitIntoSections(text string, numSections int) []string {
	sections := make([]string, numSections)
	if text == "" || numSections <= 0 {
		return sections
	}

	parts := strings.Split(text, "\n\n")

	for i := 0; i < numSections && i < len(parts); i++ {
		sections[i] = strings.TrimSpace(parts[i])
	}

	return sections
}

// Helper function to parse meal sections
func parseMealSections(text string) map[string]string {
	return map[string]string{
		"breakfast": extractSection(text, "BREAKFAST"),
		"lunch":     extractSection(text, "LUNCH"),
		"dinner":    extractSection(text, "DINNER"),
		"snacks":    extractSection(text, "SNACKS"),
	}
}

// Helper function to extract section from text
func extractSection(text, keyword string) string {
	if text == "" {
		return "No information available."
	}

	textUpper := strings.ToUpper(text)
	keywordUpper := strings.ToUpper(keyword)

	start := strings.Index(textUpper, keywordUpper)
	if start == -1 {
		return "No information available."
	}

	// Move past the keyword
	content := text[start+len(keywordUpper):]

	// Stop at the next section header
	for _, next := range []string{"BREAKFAST", "LUNCH", "DINNER", "SNACKS"} {
		if next == keywordUpper {
			continue
		}
		if idx := strings.Index(strings.ToUpper(content), next); idx != -1 {
			content = content[:idx]
			break
		}
	}

	clean := strings.TrimSpace(strings.Trim(content, ":-"))
	if clean == "" {
		return "No information available."
	}

	return clean
}

// Helper function to check for keywords
func containsKeywords(text string, keywords []string) bool {
	if text == "" || len(keywords) == 0 {
		return false
	}

	textLower := strings.ToLower(text)

	for _, keyword := range keywords {
		if keyword == "" {
			continue
		}
		if strings.Contains(textLower, strings.ToLower(keyword)) {
			return true
		}
	}

	return false
}
