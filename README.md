# Online Forum REST Server

## Project Description

This project involves the development of a **REST server** for an **imaginary online forum** where users can create accounts and participate in discussion threads. The server includes several key features such as user authentication, form validation, and a database designed using **SQLite** to store user data and forum threads.

The project also includes various security measures such as **security headers** for protection, **request logging** for tracking server activity, and **template caching** to optimize page load times.

### Key Features:
- **User Authentication:** Allows users to create accounts, log in, and manage their profiles.
- **Discussion Threads:** Users can participate in discussion threads, post replies, and view the content.
- **SQLite Database:** The forum data (user accounts, threads, and posts) is stored in an SQLite database.
- **Form Validation:** Ensures that input data is valid and secure before being stored in the database.
- **Security Headers:** Implements headers to secure HTTP requests and protect user data.
- **Template Caching:** Improves performance by caching HTML templates.

### Technologies Used:
- **Golang** (Backend Development)
- **SQLite** (Database)
- **HTML** (Frontend)
- **RESTful API** (Design and Development)
- **SQL** (Database Queries)
- **User Authentication** (Sessions and Token Management)
- **Database Normalization** (Efficient Data Structure)
- **Form Validation** (Data Integrity and Security)
- **Security Headers** (Protecting Data)
- **Request Logging** (Auditing)
  
## Endpoints

### Account Routes
- **GET `/account/create`**: Displays the form to create a new account.
- **POST `/account/create`**: Submits the form to create a new user account.
- **GET `/account/view/{id}`**: Views the details of a user account (protected route).
- **GET `/account/login`**: Displays the login form.
- **POST `/account/login`**: Submits the form to log a user into their account.
- **POST `/account/logout`**: Logs the user out of their account (protected route).

### Thread Routes
- **GET `/thread/create`**: Displays the form to create a new discussion thread (protected route).
- **POST `/thread/create`**: Submits the form to create a new thread (protected route).
- **GET `/thread/view/{id}`**: Views the details of a specific thread.

### Message Routes
- **GET `/thread/view/{id}/message/create`**: Displays the form to create a new message within a thread (protected route).
- **POST `/thread/view/{id}/message/create`**: Submits the form to post a new message within a thread (protected route).

## Database Structure

The project uses an **SQLite** database to store user accounts, threads, and posts. It follows the principles of **database normalization** to ensure data consistency and minimize redundancy.

## Security

- **Security Headers:** Implements key security headers (such as Content Security Policy and X-Content-Type-Options) to protect user data.
- **User Authentication:** Uses session management and token authentication to securely handle user logins.
- **SQL Injection Protection:** Ensures that all database queries are parameterized to prevent SQL injection attacks.

## Caching

The application caches HTML templates to reduce server load and speed up page rendering for returning users.