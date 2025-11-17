📦 Expense Tracker API (Go)

A clean, modular, production-style REST API for tracking expenses — built with Go, fully tested, and structured using a modern project layout (cmd/, internal/, handlers, storage, models, etc.).

Supports:
	•	➕ Add expenses
	•	📄 List all expenses
	•	🔍 Get expense by ID
	•	❌ Delete expense
	•	📊 Summary totals (overall + per category)
	•	💾 JSON file persistence
	•	🧪 Unit tests for models, handlers, and storage

⸻
### 📌 API Endpoints

| Feature              | Endpoint         | Method |
|----------------------|------------------|--------|
| Create an expense    | `/expenses`      | POST   |
| List all expenses    | `/expenses`      | GET    |
| Get specific expense | `/expenses/{id}` | GET    |
| Delete an expense    | `/expenses/{id}` | DELETE |
| Summary totals       | `/summary`       | GET    |

🛠 Tech Stack
	•	Go (Golang) – net/http, json, sync
	•	Modular internal architecture
	•	TDD-ready structure
	•	JSON file persistence

🔮 Future Enhancements
	•	PUT /expenses/{id} (update)
	•	Pagination & filtering
	•	Monthly breakdown endpoint
	•	SQLite or Postgres storage backend
	•	Dockerfile + Compose
	•	JWT auth (user accounts)
	•	React or Next.js frontend dashboard
	•	gRPC version

  👤 Author

Daud Abdi
Backend Developer (Go)
📍 London, UK
🔗 LinkedIn: https://www.linkedin.com/in/daudabdi0506
💻 GitHub: https://github.com/Daudsaid

