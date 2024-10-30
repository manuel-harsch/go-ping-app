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

### 3. Set up project directory structure
Create the necessary directories and files for the project structure.

```
mkdir -p cmd internal assets/static
touch cmd/main.go internal/webserver.go internal/config.go internal/ping.go internal/db.go internal/logging.go assets/static/index.html

```

**Note**: on Windows you can use 
> wsl mkdir ...
> wsl touch ...

(WSL must be installed)

- ***cmd/main.go:*** Entry point of the Go application.
- ***internal/webserver.go:*** Contains the HTTP server logic and routes.
- ***internal/config.go:** Handles loading and saving configuration.
- ***internal/ping.go:*** Implements the cyclic ping logic.
- ***internal/db.go:*** Manages interactions with SQLite.
- ***internal/logging.go:*** Handles logging.
- ***assets/static/index.html:*** HTML file for your GUI.

### 4. Install dependencies (Optional libraries)
If you're planning to use libraries like Gin, go-sqlite3, or go-ping/ping, you can install them at this step.

To install Gin for HTTP routing:
> go get github.com/gin-gonic/gin

To install SQLite driver:
> go get github.com/mattn/go-sqlite3

To install go-ping/ping for ICMP ping:
> go get github.com/go-ping/ping

### 5. Verify setup
You can now add a basic main.go file to ensure everything works before proceeding with further coding.

Edit the cmd/main.go file with a basic "Hello, World" web server:
```go
package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    log.Println("Server running at http://localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

```

### 6. Run the application
Now, run the Go program to ensure the basic setup is working.
> go run cmd/main.go

Open your browser and navigate to http://localhost:8080/ to see "Hello, World!" displayed.


## Implement configuration loading and saving (JSON).
We'll create a configuration file that stores the user settings (e.g., ping target, cycle time) and provide functions to load and save these settings.

### 1. Create a Configuration Struct
First, we need a struct that represents the configuration data. Let’s assume we are storing the following settings:

- ***Host***: The DNS name or IP address to ping.
- ***CycleTime***: Time interval between pings (in seconds).
- ***PingTimeout***: Maximum timeout for each ping (in milliseconds).

Here’s an example struct for the configuration:
```Go
// internal/config.go
package internal

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "log"
)

// Config represents the structure of the JSON configuration file
type Config struct {
    Host        string `json:"host"`
    CycleTime   int    `json:"cycle_time_seconds"`  // Ping cycle time in seconds
    PingTimeout int    `json:"ping_timeout"` // Ping timeout in milliseconds
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var config Config
    err = json.Unmarshal(bytes, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

// SaveConfig saves the configuration to a JSON file
func SaveConfig(config *Config, filePath string) error {
    bytes, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(filePath, bytes, 0644)
    if err != nil {
        return err
    }

    return nil
}
```
Explanation:
- ***Config*** struct defines the properties we want to store: Host, CycleTime, and PingTimeout.
- ***LoadConfig(filePath string)*** opens the configuration file, reads its contents, and unmarshals it into a Config struct.
- ***SaveConfig(config *Config, filePath string)*** takes a Config struct and writes its contents to a JSON file.

### 2. Create a Default Configuration File
You can also provide a default configuration in case the file doesn’t exist yet. Here's a simple utility function to create a default configuration:
```Go
// internal/config.go
func DefaultConfig() *Config {
    return &Config{
        Host:        "8.8.8.8",  // Default to Google Public DNS
        CycleTime:   5,          // 5 seconds between pings
        PingTimeout: 1000,       // 1000 ms (1 second) ping timeout
    }
}
```

### 3. Use the configuration in the main function
Now, let's modify the cmd/main.go to load the configuration when the application starts and save it when the user makes changes (e.g., through the GUI).

Here’s an example of how to load the configuration and handle the scenario where the configuration file does not exist yet:

```Go
// cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/manuel-harsch/go-ping-app/internal"
)

const configFilePath = "config.json"

func main() {
	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// File does not exist, create a default configuration
		log.Println("Config file not found, creating default config.")
		defaultConfig := internal.DefaultConfig() // Call DefaultConfig from internal package
		if err := internal.SaveConfig(defaultConfig, configFilePath); err != nil {
			log.Fatalf("Failed to create default config: %v", err)
		}
	}

	// Load the configuration
	cfg, err := internal.LoadConfig(configFilePath) // Call LoadConfig from internal package
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Printf("Loaded Config: Host=%s, CycleTime=%d ms, PingTimeout=%d ms", cfg.Host, cfg.CycleTime, cfg.PingTimeout)

	// Now you can use the loaded configuration for further processing
}
```

### 4. JSON configuration file format
The config.json file will look something like this when saved:
```JSON
{
  "host": "8.8.8.8",
  "cycle_time_seconds": 5,
  "ping_timeout": 1000
}
```

### 5. Testing the setup
To test the setup:

- Delete the config.json file (if it exists) to see if the default configuration is created.
- Run the application again

> go run cmd/main.go

Modify the ***config.json*** file and reload the application to confirm that it reads the updated configuration.


## Implement the Webserver with Routes for the GUI and API
We will set up:

- A web server using Gin (or the standard net/http package).
- API routes for:
    - Getting and updating the configuration (via GUI).
    - Starting/stopping the ping process.
- Serving the static HTML files (the GUI).

### Step-by-Step Implementation:

#### 1. Install Gin:

If you haven’t already, install Gin for building the web server.
> go get github.com/gin-gonic/gin

#### 2. Modify main.go to Set Up the Webserver:

We’ll now update the ***cmd/main.go*** file to:
- Serve the GUI (an HTML file).
- Provide API routes to fetch/update the configuration.
- Provide API routes for starting/stopping the ping process.

Update main.go as follows:
```Go
// cmd/main.go
package main

import (
	"log"
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manuel-harsch/go-ping-app/internal"
)

const configFilePath = "config.json"

func main() {
	// Check if the configuration file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// File does not exist, create a default configuration
		log.Println("Config file not found, creating default config.")
		defaultConfig := internal.DefaultConfig() // Call DefaultConfig from internal package
		if err := internal.SaveConfig(defaultConfig, configFilePath); err != nil {
			log.Fatalf("Failed to create default config: %v", err)
		}
	}

	// Load the configuration
	cfg, err := internal.LoadConfig(configFilePath) // Call LoadConfig from internal package
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up the Gin web server
	router := gin.Default()

	// Serve static files (the GUI)
	router.Static("/static", "./assets/static")

	// API: Fetch current configuration
	router.GET("/api/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, cfg)
	})

	// API: Update the configuration
	router.POST("/api/config", func(c *gin.Context) {
		var newConfig internal.Config
		if err := c.BindJSON(&newConfig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration data"})
			return
		}

		// Save the new configuration to the file
		if err := internal.SaveConfig(&newConfig, configFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
			return
		}

		// Update the in-memory config (so changes apply without restart)
		cfg = &newConfig
		c.JSON(http.StatusOK, cfg)
	})

	// API: Start ping process
	router.POST("/api/ping/start", func(c *gin.Context) {
		// Logic to start the ping process would go here
		c.JSON(http.StatusOK, gin.H{"message": "Ping started"})
	})

	// API: Stop ping process
	router.POST("/api/ping/stop", func(c *gin.Context) {
		// Logic to stop the ping process would go here
		c.JSON(http.StatusOK, gin.H{"message": "Ping stopped"})
	})

	// Start the web server
	router.Run(":8080")
}

```

**Explanation of the Code:**

**1. Serving Static Files:**

- ```router.Static("/static", "./assets/static")```: 
  This serves static files from the ```./assets/static``` directory. You can place your frontend (HTML, CSS, JS) here.

**2. API Endpoints:**

- ```GET /api/config```: Fetches the current configuration from memory (loaded from the JSON file).
- ```POST /api/config```: Accepts new configuration data, updates the JSON file, and reloads the configuration into memory.
- ```POST /api/ping/start```: An endpoint to start the ping process (we will implement the ping logic later).
- ```POST /api/ping/stop```: An endpoint to stop the ping process.

**3. Gin's JSON Handling:**

- ```c.BindJSON(&newConfig)```: Binds the incoming JSON data (from the GUI/API request) to a Config struct.
- ```c.JSON```: Sends JSON responses back to the client.

**4. Starting the Server:**

- ```router.Run(":8080")```: Starts the server on port 8080.

---

#### 3. Create a Simple HTML GUI:

Let’s create a simple HTML file to serve as the frontend GUI for configuration and ping control. 

Place the following file in ```assets/static/index.html```.

```HTML
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ping App</title>
    <script>
        async function loadConfig() {
            const response = await fetch('/api/config');
            const config = await response.json();
            document.getElementById('host').value = config.host;
            document.getElementById('cycle_time').value = config.cycle_time_milliseconds;
            document.getElementById('ping_timeout').value = config.ping_timeout_milliseconds;
        }

        async function saveConfig() {
            const host = document.getElementById('host').value;
            const cycleTime = document.getElementById('cycle_time').value;
            const pingTimeout = document.getElementById('ping_timeout').value;
            const config = {
                host: host,
                cycle_time_milliseconds: parseInt(cycleTime),
                ping_timeout_milliseconds: parseInt(pingTimeout)
            };

            const response = await fetch('/api/config', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(config)
            });

            if (response.ok) {
                alert("Configuration saved successfully!");
            } else {
                alert("Failed to save configuration.");
            }
        }

        async function startPing() {
            const response = await fetch('/api/ping/start', {
                method: 'POST'
            });

            if (response.ok) {
                alert("Ping started!");
            } else {
                alert("Failed to start ping.");
            }
        }

        async function stopPing() {
            const response = await fetch('/api/ping/stop', {
                method: 'POST'
            });

            if (response.ok) {
                alert("Ping stopped!");
            } else {
                alert("Failed to stop ping.");
            }
        }

        window.onload = loadConfig;
    </script>
</head>
<body>
    <h1>Ping App Configuration</h1>
    <label>Host: <input type="text" id="host" /></label><br>
    <label>Cycle Time (ms): <input type="number" id="cycle_time" /></label><br>
    <label>Ping Timeout (ms): <input type="number" id="ping_timeout" /></label><br>
    <button onclick="saveConfig()">Save Configuration</button><br><br>

    <h2>Ping Control</h2>
    <button onclick="startPing()">Start Ping</button>
    <button onclick="stopPing()">Stop Ping</button>
</body>
</html>

```
**Explanation:**

- **JavaScript for API Interaction:**
  - Fetches the configuration from the ```/api/config``` endpoint and fills the HTML form with current values.
  - Allows the user to update the configuration and save it by sending a ```POST``` request to ```/api/config```.
  - Provides buttons to start/stop the ping process, interacting with ```/api/ping/start``` and ```/api/ping/stop```.

---
#### 4. Testing the Application:

1. Run your Go application:

	> go run cmd/main.go

2. Open your browser and go to ```http://localhost:8080/static/index.html``` to see the simple GUI.

3. You should be able to:
   - View the current configuration.
   - Update and save the configuration.
   - Use the start/stop buttons (we will implement the actual ping logic in the next step).
