# 🩺 DiabetesAI Advisor

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Genkit](https://img.shields.io/badge/Firebase-Genkit-orange)
![Gemini](https://img.shields.io/badge/Google-Gemini_AI-blue?logo=google)
![Cloud Run](https://img.shields.io/badge/Google%20Cloud-Run-4285F4?logo=googlecloud&logoColor=white)
![Healthcare](https://img.shields.io/badge/Healthcare-AI-green)
![CI/CD](https://github.com/YOUR_USERNAME/diabeticai-advisor/actions/workflows/deploy.yml/badge.svg)

An AI-powered diabetes management assistant built with **Go** and **Firebase Genkit**, powered by Google's **Gemini 2.5 Flash AI Model**.

DiabetesAI Advisor provides safe, educational guidance to help people living with diabetes better understand their health data and daily choices.

---

## 🚀 What It Does

DiabetesAI Advisor supports diabetes self-management through five AI-powered features:

1. **Blood Sugar Interpreter**  
   Understand blood glucose readings and get contextual guidance.

2. **Meal Planner**  
   Generate diabetes-friendly meal suggestions based on dietary needs.

3. **Symptom Checker**  
   Assess symptoms and get guidance on when to seek medical attention.

4. **Exercise Advisor**  
   Receive safe exercise recommendations tailored to fitness level and blood sugar.

5. **Medication Information**  
   Learn about common diabetes medications (educational use only).

> ⚠️ This tool does **not diagnose** or replace professional medical care.

---

## 🎯 Why I Built This

This project is focusing on strengthening primary healthcare in Africa through digital innovation.

Diabetes is a growing public health challenge, and responsible AI can help improve access to reliable health information—especially in resource-constrained settings.

---

## 🧠 How It Works

- Built using **Firebase Genkit** for agentic AI workflows
- Uses **Google Gemini model (Gemini 2.5 Flash)** for reasoning and text generation
- Exposes structured HTTP endpoints for each health capability
- Designed for **safe, non-diagnostic, human-readable outputs**

---

## 🛠 Technology Stack

- **Language:** Go (Golang)
- **AI Framework:** Firebase Genkit (Agentic AI orchestration)
- **AI Models:** Google Gemini (Google AI)
- **Server:** Genkit HTTP server
- **Deployment:** Google Cloud Run
- **CI/CD:** GitHub Actions

---

## ☁️ Live Demo

🔗 **Cloud Run URL:**  
https://YOUR_CLOUD_RUN_URL



---

## ⚙️ Getting Started

### Prerequisites

- Go 1.21 or higher
- Google Gemini API key
- VS Code (recommended)

---

### Installation

1. Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/diabeticai-advisor.git


cd diabeticai-advisor






Install dependencies:


go mod tidy






Create a .env file:

// Add Google's AI API Key


GEMINI_API_KEY=your_api_key_here






Load environment variables:

Linux / macOS


export $(cat .env | xargs)






Windows (PowerShell)



Get-Content .env | ForEach-Object {
  $var = $_.Split('=')
  [Environment]::SetEnvironmentVariable($var[0], $var[1])
}






Run the app:


go run main.go





The server starts at:



http://localhost:3400







🧪 How to Use (API Examples)


Blood Sugar Interpretation

curl -X POST http://localhost:3400/bloodSugar \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "reading": 145,
      "meal_timing": "after_meal",
      "meal_type": "lunch"
    }
  }'





Meal Plan

curl -X POST http://localhost:3400/mealPlan \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "diet_type": "vegetarian",
      "allergies": "none",
      "calorie_limit": 1800
    }
  }'






Symptom Checker

curl -X POST http://localhost:3400/symptoms \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "symptoms": "feeling dizzy and tired",
      "duration": "2 hours",
      "current_meds": "metformin"
    }
  }'






Exercise Advisor

curl -X POST http://localhost:3400/exercise \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "fitness_level": "beginner",
      "time_available": 30,
      "current_bg": 120,
      "preferred_type": "walking"
    }
  }'






Medication Information

curl -X POST http://localhost:3400/medication \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "medication_name": "metformin",
      "purpose": "timing"
    }
  }'







🔌 API Endpoints
Endpoint	Method	Description
/bloodSugar	POST	Interpret blood glucose readings
/mealPlan	POST	Generate diabetes-friendly meal plans
/symptoms	POST	Symptom assessment and guidance
/exercise	POST	Exercise recommendations
/medication	POST	Medication information









🔐 Security & Privacy

No user data is stored

Stateless, privacy-first design

Suitable for educational and prototype use









⚠️ Medical Disclaimer

This project is for educational purposes only.

Does not provide medical diagnosis or treatment

Always consult a qualified healthcare professional

Do not change medications without medical advice

In emergencies, contact local emergency services








🚧 Future Enhancements

Firestore integration for user history

Authentication & user profiles

Multi-language support

Blood sugar trends & analytics

Mobile / web frontend

Integration with glucose monitoring devices








🤝 Contributing
Contributions are welcome:

Bug reports

Feature suggestions

Pull requests








📄 License
This project is licensed under the MIT License.
See the full license here:
👉 https://opensource.org/licenses/MIT

Built with ❤️ using Go, Firebase Genkit, and Google Gemini