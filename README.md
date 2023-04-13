# **Replace Tokens**

## **Overview**

Replace tokens is a command line tool that helps to replace variable tokens in configuration files.\
The main goal of this tool is to help the deployment process by defining configuration variables in each deployment stage or environment.

## **Installation**

Users can run the Replace Tokens Tool in one of the following ways:

### Binaries

Binary files can be found in the release section. Users can download the binaries and execute them from the command line or terminal. The available platforms are the following.

* [Linux x86-amd64](https://github.com/enaldo1709/replace-tokens/releases/download/replacetokens/replacetokens-linux)
* [Windows x86-amd64](https://github.com/enaldo1709/replace-tokens/releases/download/replacetokens/replacetokens-windows.exe)
* [Apple Mac OS Intel Chips x86-amd64](https://github.com/enaldo1709/replace-tokens/releases/download/replacetokens/replacetokens-macos)

### Build from source

The next steps are to build the tool from the source code:

1. #### Dependencies

    Some dependencies need to be satisfied to build the tool from the source code:

    * [Go v1.20](https://go.dev/)
    * [make](https://www.gnu.org/software/make/)

2. #### Download Sources

    Source code can be downloaded from the following sources:

    * [Release Page - zip](https://github.com/enaldo1709/replace-tokens/archive/refs/tags/replacetokens.zip)
    * [Release Page - tar](https://github.com/enaldo1709/replace-tokens/archive/refs/tags/replacetokens.tar.gz)
    * Clone git repository:

        ```bash
        git clone https://github.com/enaldo1709/replace-tokens.git
        ```

3. #### Building

    Build for Unix os (Linux and Mac Os) generate an executable file named *replacetokens* in the root project folder:

    ```bash
    make build
    ```

    Build for Windows, generate an executable file named *replacetokens.exe* in the root project folder:

    ```bash
    make build-windows
    ```

4. #### Install

    Install in $GOBIN folder:

    ```bash
    make install
    ```

    Install in linux root /bin folder (requires root):

    ```bash
    make install-to-system
    ```

## **Usage**

### Command Line

```bash
replacetokens [PREFIX] [SUFFIX] [TOKENS FILE] [FILE TO REPLACE] [OUTPUT]
```

| argument      | description| Required |
|---------------|------------|----------|
| PREFIX        | Prefix used to denote a token in the file to replace. | required |
| SUFFIX        | Suffix used to denote a token in the file to replace. | required |
| TOKENS FILE   | Path to the file that contains the values of the tokens to be replaced. Must be key-value in YAML format with just one hierarchical level. Eg. **TOKEN: value** | required |
| FILE TO REPLACE | Path to the file that contains the tokens which might be replaced by the values in the TOKENS FILE | required |
| OUTPUT | Path to the file where will be wrote the FILE TO REPLACE with all tokens replaced. If OUTPUT is not provided, the output file will be paced in the same location oh the FILE TO REPLACE with the same name adding a flag at the end of the filename with value of "-replaced" | optional |

### Usage Example

This is an example of how to use replace tokens tool, for this example, the application binary executable is located in a PATH folder (in this case in the $GOBIN folder) so is accessible through the command line. See the [Installation](#installation) section.\
Imagine an application that uses a properties configuration file. There are two environments defined for this application: beta and prod.\
The project structure is as follows:

```bash
 .
├──  go.mod
├──  src
│   ├──  main.go
│   ├──  config.properties
├──  config
│   ├──  beta.yaml
│   └──  prod.yaml
```

Think that the application needs 3 properties: log level, database user and database password.

* *config.properties*

    ```bash
    log.level={LOG_LEVEL}
    database.user={DB_USER}
    database.password={DB_PASSWORD}
    ```

So the config values present in token files (beta.yaml and prod.yaml) are:

* *beta.yaml*

    ```yaml
    LOG_LEVEL: debug
    DB_USER: beta-user
    DB_PASSWORD: beta-password
    ```

* *prod.yaml*

    ```yaml
    LOG_LEVEL: info
    DB_USER: prod-user
    DB_PASSWORD: prod-password
    ```

Using the following commands **replace tokens tool** can replace all values in the *config.properties* file:

```bash
replacetokens "{" "}" config/beta.yaml src/config.properties src/config-beta.properties
replacetokens "{" "}" config/prod.yaml src/config.properties src/config-prod.properties
```

The new project structure is:

```bash
 .
├──  go.mod
├──  src
│   ├──  main.go
│   ├──  config.properties
│   ├──  config-beta.properties
│   ├──  config-prod.properties
├──  config
│   ├──  beta.yaml
│   └──  prod.yaml
```

The output files are:

```bash
project$ cat src/config-beta.properties
log.level=debug
database.user=beta-user
database.password=beta-password

project$ cat config-prod.properties
log.level=info
database.user=prod-user
database.password=prod-password
```

## **Author**

**Enaldo Narvaez Yepes**:\
Software Developer.

* [Github](https://github.com/enaldo1709)
* [LinkedIn](https://www.linkedin.com/in/enaldo-narv%C3%A1ez-3386b6148/)

## **Thanks**

Thanks to all open source and Go community for the effort of building all go platforms and tools.
