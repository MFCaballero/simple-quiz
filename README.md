# Simple Quiz Application

This README provides instructions for setting up and using the Quiz application, which consists of a backend server and a command-line interface (CLI) for interacting with quizzes.

## Prerequisites

Before running the Quiz application, ensure you have the following prerequisites installed:

- Go programming language: [https://golang.org/](https://golang.org/)

## Getting Started

- Clone the repository:

```bash
git clone <repository-url>
cd <repository-directory>
```

- Start the backend server and build the CLI binary:

```bash
make start-quiz
```
This command will run the backend server on port 8080 and build the CLI binary. The backend service will be running on the terminal used.
## Using the CLI

Open a new terminal and navigate to the project directory.

### Help Command

Use the help command to view available commands:
```bash
./quiz help
```

### Login Command
Login to the quiz app

```bash
./quiz login [flags]
```
#### Flags
-u, --userName string: Your user name

#### Example
```bash
./quiz login -u Maria
```

### Question Command
Interact with quiz questions

```bash
./quiz question [command]
```
#### Available Commands:

#### List all questions
```bash
./quiz question list
```
#### Get a particular question with it's question number
```bash
./quiz question get [flags]
```
<p>Flags</p>
-h, --help   help for get  <br>
-n, --questionNumber string   Question number
<br></br>
<p>Example:</p>

```bash
./quiz question get -n 1
```

### Answer Command
Interact with quiz
```bash
./quiz answer [command]
```
#### Available Commands:

#### Answer a quiz question
You can answer the same question multiple times, the quiz will only save the last posted answer
```bash
./quiz answer post [flags]
```
<p>Flags</p>
-h, --help   help for get  <br>
-o, --option string Option letter <br>
-q, --question string Question number
<br></br>
<p>Example:</p>

```bash
./quiz answer post -q 1 -o A
```

#### Get answered questions
```bash
./quiz answer list
```

#### Finish the quiz
This command will post all your answers and you'll no longer be able to edit an answer. It's only available if all the questions are completed
```bash
./quiz answer finish
```

#### Get your Score results

```bash
./quiz answer score
```

### Logout Command
Logout from the quiz app

```bash
./quiz logout
```