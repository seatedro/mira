# Mirame

Mirame is an interpreted programming language written in Go. It is designed to follow the principles outlined in Thorsten Ball's book "Writing an Interpreter in Go" and serves as a practical example of how to write an interpreter from scratch using Go.
Mira is the repository that holds the interpreter code written in Go. [Mirame](https://github.com/seatedro/mirame) Is the repository for the interpreter written in Rust.

## Getting Started

### Prerequisites

- Go (version 1.16 or higher)

### Installing

To install Mira, clone this repository into your local machine:

```
git clone git@github.com:seatedro/mira.git
```

(If you're not using SSH, use it)

### Running

Once you have cloned the repository, navigate to the root directory of the project and run the following command:

```
go run main.go
```

This will start the Mira interpreter and you can begin executing commands.

### Usage

Mira currently supports the following commands:

- `let <variable name> = <value>;`: assigns a value to a variable
- `<variable name>;`: retrieves the value of a variable
- `<expression>;`: evaluates an expression

#### Example

```
>> let x = 5;
>> let y = 10;
>> x + y;
15
```

## Contributing

Contributions to Mira are welcome. If you have any issues or feature requests, please submit them via GitHub issues.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
