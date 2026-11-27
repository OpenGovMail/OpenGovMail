# 📊 Load Testing Suite

This directory contains the load testing suite for the OpenGovMail Mail Server using [Locust](https://locust.io/). The tests simulate real-world email traffic patterns including SMTP sending and IMAP operations.

## 📁 Files Overview

- **`locustfile.py`** - Main Locust configuration file that orchestrates all load tests
- **`smtp_tester.py`** - SMTP load testing scenarios (sending emails)
- **`imap_tester.py`** - IMAP load testing scenarios (reading, managing emails)
- **`config.py`** - Email server configuration and settings
- **`user_manager.py`** - User management utilities for test users
- **`data_generator.py`** - Generates realistic test data (emails, attachments)
- **`requirements.txt`** - Python dependencies
- **`test_data/`** - Directory containing test users and attachments

## 🚀 Quick Start

### Prerequisites

1. Python 3.8 or higher
2. OpenGovMail Mail Server running and accessible

### Step 1: Create Test Users

Before you run any load tests, create the test users as [Add users](../../docs/how-to/add-users.md) describes. Then list them in `test_data/users.csv`:

```csv
username,email,password
testuser001,testuser001@example.com,<password>
```

If the file doesn't exist, the tests generate one with sample users that aren't on the server, and every login fails.

### Step 2: Installation

1. Create and activate a virtual environment (recommended):

```bash
python3 -m venv .venv
source .venv/bin/activate  # On macOS/Linux
# or
.venv\Scripts\activate  # On Windows
```

2. Install dependencies:

```bash
pip install -r requirements.txt
```

### Configuration

Set the mail domain as an environment variable or create a `.env` file:

```bash
export MAIL_DOMAIN=your-domain.com
```

Or create a `.env` file in the `test/load` directory:

```
MAIL_DOMAIN=your-domain.com
```

## 📋 Running Load Tests

### Interactive Mode (Web UI)

Run Locust with the web interface to manually control the load test:

```bash
locust -f locustfile.py --host https://your-domain.com
```

Then open your browser to `http://localhost:8089` and configure:
- Number of users
- Spawn rate (users per second)
- Host URL

### Headless Mode (Automated)

Run Locust in headless mode with predefined parameters:

```bash
locust \
  -f locustfile.py \
  --headless \
  -u 10 \
  -r 2 \
  --run-time 60s \
  --host https://your-domain.com \
  --html report.html \
  --csv results \
  --logfile locust.log \
  --loglevel INFO
```

#### Parameters Explained:

- `-f locustfile.py` - Locust file to use
- `--headless` - Run without web UI
- `-u 10` - Number of concurrent users to simulate
- `-r 2` - Spawn rate (users spawned per second)
- `--run-time 60s` - Test duration (60 seconds)
- `--host https://your-domain.com` - Target mail server
- `--html report.html` - Generate HTML report
- `--csv results` - Generate CSV results (creates results_stats.csv, results_failures.csv, results_stats_history.csv)
- `--logfile locust.log` - Log file location
- `--loglevel INFO` - Logging level

### Custom Test Scenarios

#### Light Load Test (5 users, 30 seconds)

```bash
locust \
  -f locustfile.py \
  --headless \
  -u 5 \
  -r 1 \
  --run-time 30s \
  --host https://your-domain.com \
  --html report-light.html
```

#### Heavy Load Test (50 users, 5 minutes)

```bash
locust \
  -f locustfile.py \
  --headless \
  -u 50 \
  -r 5 \
  --run-time 5m \
  --host https://your-domain.com \
  --html report-heavy.html \
  --csv results-heavy
```

#### SMTP Only Test

```bash
locust \
  -f locustfile.py \
  --headless \
  -u 10 \
  -r 2 \
  --run-time 60s \
  --host https://your-domain.com \
  --tags smtp
```

#### IMAP Only Test

```bash
locust \
  -f locustfile.py \
  --headless \
  -u 10 \
  -r 2 \
  --run-time 60s \
  --host https://your-domain.com \
  --tags imap
```

## 📊 Understanding Results

### HTML Report

The HTML report (`report.html`) includes:
- Request statistics (RPS, response times, failures)
- Charts showing performance over time
- Percentile response times (50th, 66th, 75th, 80th, 90th, 95th, 98th, 99th, 100th)
- Total request count and failure rate

### CSV Results

Three CSV files are generated:
- `results_stats.csv` - Aggregated statistics per request type
- `results_failures.csv` - Details of all failures
- `results_stats_history.csv` - Time-series data of statistics

### Log File

The `locust.log` file contains detailed logs of the test execution, useful for debugging failures.

## 📚 Resources

- [Locust Documentation](https://docs.locust.io/)
- [Writing Locust Tests](https://docs.locust.io/en/stable/writing-a-locustfile.html)
- [Locust Configuration](https://docs.locust.io/en/stable/configuration.html)

## 🆘 Support

For issues or questions:
1. Check the log files (`locust.log`)
2. Review the test/load test results
3. Consult the OpenGovMail Mail documentation
4. Open an issue in the repository
