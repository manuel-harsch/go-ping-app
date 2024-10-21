# go-ping-app - Simple PING stats app

# --- Planning ---

### Components Breakdown

#### 1. GUI - served Web App
- You can use Go's net/http package to serve a web app and use a frontend framework (like HTML/CSS/JavaScript) for the user interface. To simplify things, you can use a Go framework like Gin or Fiber.
- For the frontend, frameworks like Vue.js, React, or even simple Vanilla JS for handling user interactions (such as ping configuration, start/stop buttons, and displaying data).

#### 2. Configuration (JSON File)
- Use encoding/json to read and write configurations to a JSON file. The configuration should store host details, cycle time, etc.
- Whenever a user modifies settings via the GUI, these updates should be saved to the configuration file.
  
#### 3. SQLite / Time-Series database
- SQLite is a great fit for a lightweight, serverless database.
- If you want a time-series database, you can consider InfluxDB which is designed for time-series data. However, SQLite can also handle time-series data by including a timestamp column.
- For Go, the mattn/go-sqlite3 driver works well with SQLite.

#### 4. Logging
- Use Go's built-in log package for logging events such as pings, errors, and configuration changes.
- You can also use libraries like Logrus or Zap for structured and leveled logging.


### Features breakdown

#### 1. Single executable file

- You can compile the entire application into a single binary using Go's go build command.
- Use Go modules for dependency management.

#### 2. Configurable host, DNS name, or IP to ping (GUI)
- Add input fields in the GUI to allow users to configure the host/IP to ping.
- Save these settings to the JSON configuration file using Go’s file I/O methods.

#### 3. Configurable cycle time (GUI)
- Provide an input for the ping interval (in seconds/minutes) on the GUI.
- Save the cycle time to the configuration JSON file and update it in real-time as needed.

#### 4. Cyclic Ping of configured host
- You can use Go's ICMP package (github.com/go-ping/ping) to handle pings.
- Use Go routines to run the ping process cyclically at intervals based on the configured cycle time.
- Save the ping results to the SQLite database (with timestamps, latency, etc.).

#### 5. Save Ping results to database
- Define a table in SQLite with columns like timestamp, target_host, latency, etc.
- Insert each ping result into the SQLite database with a Go function triggered after each ping.

#### 6. Start/Stop button on GUI
- The start/stop button can toggle the cyclic ping process. You can handle this with Go routines or channels to pause or resume the pings.

#### 7. Tab to show database content (Table View)
- Use SQL queries to fetch data from the SQLite database and display it on the GUI in a tabular format.
- Allow the user to apply filters, such as date range or target host.

#### 8. Tab to show database content (Graph View)
- Use JavaScript graphing libraries like Chart.js or D3.js to visualize the ping statistics on the GUI.
- Provide date selection fields for filtering the displayed data by time range.


## Basic Go project structure

```
- /cmd
  - main.go              # Entry point
- /internal
  - webserver.go         # HTTP server and routes
  - config.go            # Configuration handling (JSON read/write)
  - ping.go              # Ping logic (ICMP, periodic pings)
  - db.go                # SQLite database handling
  - logging.go           # Logging
- /assets
  - /static              # Frontend files (HTML/CSS/JS)
- config.json            # Configuration file
- go.mod                 # Dependency management
```

## Libraries to consider

### Frontend
- HTML/CSS/JavaScript (or Vue.js/React.js)
- Chart.js/D3.js (for graphs)

### Backend
- net/http or Gin/Fiber for web serving.
- go-ping/ping for ICMP pings.
- mattn/go-sqlite3 for SQLite.
- encoding/json for configuration.
- log or Logrus/Zap for logging.
- goroutines for background ping cycles.


# --- Next steps planning ---
1. Set up a Go project with basic structure (go mod init).
2. Implement configuration loading and saving (JSON).
3. Implement the webserver with routes for the GUI and API (for configuration and ping control).
4. Implement ping logic using goroutines and SQLite to save ping results.
5. Create the frontend with options for configuration, start/stop, and displaying statistics.

# --- Implementation ---
## Set up a Go project with basic structure
### 1. Create a project directory

First, create a new directory for your Go project and navigate into it.

```
mkdir go-ping-app
cd go-ping-app
```

### 2. Initialize a Go module
Run the following command to initialize a new Go module. This creates a go.mod file that helps manage dependencies.

```
go mod init github.com/manuel-harsch/go-ping-app
```
