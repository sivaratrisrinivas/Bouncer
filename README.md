# Bouncer: Supabase Rate Limiter🚦

This project is a working demonstration of a "concurrency shaper" or "rate limiter" built in Go that sits in front of a Supabase database.

## 🚀 Live Demo

You can try out the live application here: **[https://bouncer-2r9n.onrender.com](https://bouncer-2r9n.onrender.com)**

It solves the "noisy neighbor" problem. Imagine living in an apartment with shared Wi-Fi. If one person starts downloading huge files, it slows down the internet for everyone else. This project acts like a smart router that gives each user a fair "bandwidth budget," ensuring the system remains fast and reliable for all tenants.

## How It Works

The project is made of three parts that work together:

1.  **Frontend Website (`index.html`, `script.js`)**: A simple webpage with buttons that the user clicks. This is the user interface.
2.  **Go Backend (`main.go`)**: The "brains" of the operation. This is a smart web server that receives requests from the frontend. It checks who the user is (a "free" or "pro" tier user) and uses a **Token Bucket** algorithm to decide if they are within their usage limit. If they are, it fetches data from Supabase. If they're not, it tells them to wait.
3.  **Supabase Database**: A real PostgreSQL database hosted on Supabase that stores the sample `products` data.

The flow is simple: The **Frontend** talks to the **Go Backend**, which then talks to the **Supabase Database**.

---

## Features

* **Tenant-Based Rate Limiting**: Each user "tier" gets its own set of rules.
* **"Free Tier"**: Has a stricter limit, allowing fewer requests in a short period.
* **"Pro Tier"**: Has a more generous limit, allowing for much more traffic.
* **Real-time Feedback**: The frontend immediately shows if a request was successful or if it was rate-limited.

## 📸 Screenshot

*A preview of the Bouncer Glassmorphism UI in action.*

![Bouncer UI](https://i.imgur.com/T0b6G7N.png)

---

## Technology Stack

* **Backend**: Go (`net/http` for the server, `pgx` for database connection)
* **Database**: Supabase (PostgreSQL)
* **Frontend**: Plain HTML, CSS, and JavaScript

---

## How to Set Up and Run This Project

Follow these steps to get the project running on your own machine.

### 1. Prerequisites

* You must have **Go** installed on your computer.
* You need a free **Supabase** account.

### 2. Set up Supabase

1.  Create a new project in your Supabase dashboard.
2.  Go to the **SQL Editor** and run the following SQL to create the `products` table and add some data:
    ```sql
    create table products (
      id int primary key,
      name text,
      price int
    );

    insert into products (id, name, price) values
    (1, 'Laptop', 1200),
    (2, 'Mouse', 25),
    (3, 'Keyboard', 75);
    ```
3.  Go to **Settings** -> **Database** and copy your connection string. It will start with `postgres://`.

### 3. Configure the Go Backend

1.  Open the `main.go` file.
2.  Find this line:
    ```go
    connStr := "postgres://postgres:[YOUR-PASSWORD]@db.xxxx.supabase.co:5432/postgres"
    ```
3.  **Paste your own Supabase connection string** here.

### 4. Run the Application

1.  Open your computer's terminal or command prompt.
2.  Navigate to your project folder (the one with `main.go`, `index.html`, etc.).
3.  Run the following command:
    ```bash
    go get github.com/jackc/pgx/v5
    ```
4.  Now, start the server:
    ```bash
    go run main.go
    ```
5.  Open your web browser and go to this address: **`http://localhost:8080`**

---

## How to Test the Rate Limiter

Once the website is open in your browser:
* Click the **"Query as Pro User"** button. The data should appear instantly.
* Now, click the **"Query as Free User"** button repeatedly and very quickly (about 10-15 times). You will see the results box turn into an error message: **`Error: 429 - Too Many Requests`**. This shows the rate limiter is working! 