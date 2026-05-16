```markdown
# 🛒 Gadget Shop API (NoSQL Advanced Project)

**Course:** Advanced Databases (NoSQL)  
**Team:** [Your Name] & Nurlan Bekov  
**Project Type:** E-commerce Backend API with MongoDB & Go

## 📖 Project Overview
This project is a RESTful API for an online electronics store. It is built to demonstrate **advanced MongoDB patterns**, including hybrid data modeling (embedding vs. referencing), complex aggregations, and performance optimization using indexes.

### 🛠 Tech Stack
* **Language:** Go (Golang)
* **Database:** MongoDB (v6.0 via Docker)
* **Containerization:** Docker & Docker Compose
* **Admin UI:** Mongo Express / MongoDB Compass

---

## 🚀 Getting Started

### Prerequisites
* Docker & Docker Compose installed.
* Go (Golang) 1.20+ installed.

### 1. Start the Database
We use Docker to ensure a consistent environment. Run the following command to start MongoDB and the Admin UI:

```bash
docker-compose up -d

```

* **MongoDB URL:** `mongodb://localhost:27017`
* **Admin UI:** `http://localhost:8081`

### 2. Seed the Database (Generate Dummy Data)

To demonstrate aggregations and pagination, we have a script that generates **250+ realistic documents** (Users, Products, Orders).

```bash
go run seed.go

```

* *Output:* ✅ Created 150 Orders, 50 Customers, etc.

---

## 🔎 How to Check the Database

You can verify the data using three different methods:

### Option A: Web Interface (Mongo Express) - *Easiest*

We included a web-based admin panel in the Docker setup.

1. Open your browser to: **[http://localhost:8081](https://www.google.com/search?q=http://localhost:8081)**
2. **Login Credentials:**
* Username: `admin`
* Password: `pass`


3. Click on **`gadget_shop`** to view collections and documents.

### Option B: Terminal (CLI) - *Fastest*

You can check the data directly inside the container without installing anything.

```bash
# 1. Log into the database container
docker exec -it gadget_shop_db mongosh -u admin -p password123

# 2. Switch to the correct database
use gadget_shop

# 3. Check data counts
db.orders.countDocuments()  # Should be ~150
db.products.findOne()       # Check schema structure

```

### Option C: MongoDB Compass - *Best for Visuals*

If you prefer using the desktop GUI:

1. Open **MongoDB Compass**.
2. Use this **Connection String**:
```
mongodb://admin:password123@localhost:27017/?authSource=admin

```


3. Connect and browse the `gadget_shop` database.

---

## 🗄 Database Architecture (Schema Design)

We utilized a **Hybrid Data Model** to balance read performance and data consistency.

### 1. `Products` Collection

* **Pattern:** Referenced `Supplier` (Normalization) but Embedded `Specifications` (Denormalization).
* **Why?** Suppliers change rarely, but specifications are read every time a product is viewed.
* **Indexes:**
* `{ name: "text", brand: "text" }`: For search functionality.
* `{ category: 1, price: -1 }`: Compound index for sorting cheap items in categories.



### 2. `Orders` Collection

* **Pattern:** Embedded `Items`.
* **Why?** An order is a snapshot in time. Even if a product price changes later, the price *in the order* must remain the same.
* **Fields:**
* `items`: Array of objects (Snapshot of product name/price).
* `status`: Enum (`processing`, `shipped`, `delivered`).



### 3. `Customers` Collection

* **Pattern:** Embedded `Address`.
* **Why?** An address is strictly coupled to a user and rarely queried independently.

---

## ⚡ Key Features (Rubric Requirements)

### ✅ Advanced CRUD

* **Deep Updates:** Using `$inc` to decrease stock and `$set` to update statuses.
* **Transactional Logic:** (Planned) Ensuring stock is only reduced if the order is successfully created.

### 📊 Aggregations

The API includes aggregation pipelines to calculate business metrics:

1. **Sales Reports:** Grouping orders by month to show revenue trends.
2. **Top Suppliers:** Calculating which supplier generates the most profit.

### 🔍 Optimization

* **Compound Indexes** are used on `Products` to speed up filtering.
* **Projections** are used in API responses to return only necessary fields (reducing network load).

---

## 📂 Project Structure

```bash
.
├── docker-compose.yml   # Database configuration
├── seed.go              # Data generation script
├── README.md            # Documentation
├── go.mod               # Go dependencies
└── src
    ├── config           # DB Connection logic
    ├── controllers      # Request handlers
    ├── models           # Go Structs & Schema definitions
    └── routes           # API Route definitions

```

---

## 📝 License

Academic Project - SE-2434

```

```