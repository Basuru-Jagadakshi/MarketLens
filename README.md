# MarketLens
The MarketLens project is an open-source initiative designed to crawl, process and analyze labour market data in Sri Lanka. By transforming unstructured job postings into structured intelligence, this platform provides actionable insights into industry trends, skill demands and occupational shifts.

🚀 Key Functionalities
National Market Overview: Real-time metrics on vacancies, trends, and distributions by occupation, industry, and education.

Skill Intelligence: Deep-dive analysis into industry-specific skill demands and top hiring employers.

Historical Trends: Three-year longitudinal analysis to track hiring momentum.

Deep Analytics: Spatial and demographic insights into job market allocation.

🏛 Data Standards & Methodology
To ensure our labour market analytics are accurate, professional, and locally relevant, we align our data processing and classification with official Sri Lankan national standards:

Occupational Classification: We use the Sri Lanka Standard Classification of Occupations (SLSCO), based on the International Standard Classification of Occupations (ISCO-08) framework. This allows us to map vacancies across the 10 standard occupational bands effectively.

Industrial Classification: All industry-level insights are categorized according to the Sri Lanka Standard Industrial Classification (SLSIC), which is fully aligned with ISIC Rev. 4. This ensures our industry divisions are consistent with national economic reporting.

🏗 Architecture
The system is built as a modular, containerized application designed for scalability and deployment on Kubernetes.

Crawler System
Our custom crawler, powered by Crawl4AI and enriched via DeepSeek API, extracts and standardizes job data.

🛠 Tech Stack
Database: PostgreSQL

Backend: Golang (REST API)

BFF (Backend for Frontend): Next.js

Frontend: Next.js (4-view dashboard)

ORM: GORM

Data Intelligence: Crawl4AI + DeepSeek API

Deployment: Docker Compose (Rancher)

📜 License
This project is licensed under the MIT License - see the LICENSE file for details.

🚀 Running the Application

First of all you should define the following enviroment variables in the relevant folder.

1. backend
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_HOST=
DB_PORT=

POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=

2. bff
GO_BACKEND_URL=

3. crawler
DEEPSEEK_API_KEY=
BACKEND_URL=

4. frontend
NEXT_PUBLIC_API_BASE_URL=


And then restore the dump file called "market_lens_database_backup".

To deploy the stack locally or in your development environment, ensure you have Docker and Docker Compose installed. From the root directory of the project, execute the following command:

```
docker compose up --build
```
