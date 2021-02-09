package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"errors"
	"path/filepath"
	"github.com/kardianos/osext"
)

var vaultAddress = os.Getenv("VAULT_ADDR")
var filename = os.Getenv("HOME") + "/.vault_tokens"
var filevault = os.Getenv("HOME") + "/.vault"

func createFile(filename string, content string) error {
     data := []byte(content)
     err := ioutil.WriteFile(filename, data, 0664)
     if err != nil {
        return err
     }
     return nil
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func getTokens() map[string]interface{} {
     var tokens map[string]interface{}
     plan, err := ioutil.ReadFile(filename)
     if err != nil  {
        panic(err)
     }
     err2 := json.Unmarshal(plan, &tokens)
     if err2 != nil {
	panic(err2)
     }
     return tokens
}

func enable() {
    if !fileExists(filevault) {
	    bin, err := osext.Executable()
	    if err != nil {
		    panic(err)
	    }
	    content := fmt.Sprintf("token_helper = \"%s\"", filepath.ToSlash(bin))
	    er := createFile(filevault, content)
	    if er != nil {
		    panic(er)
	    }
    }
}

func get() {
    tokens := getTokens()
    if tokens[vaultAddress] != nil {
       fmt.Println(tokens[vaultAddress])
    }
    os.Exit(0)
}

func store(token string) {
    tokens := getTokens()
    tokens[vaultAddress] = token
    payload, err := json.Marshal(tokens)
    if err != nil {
       println(err)
       os.Exit(-1)
       //panic(err)
    }
    er := ioutil.WriteFile( filename, payload, 0664)
    if er != nil {
       println(err)
       os.Exit(-1)
       //panic(er)
    }
}

func erase() {
    tokens := getTokens()
    delete( tokens, vaultAddress)
    os.Exit(0)
}

func main() {
    if !fileExists(filename) {
       err := createFile(filename, "{}")
       if err != nil {
          println(err)
	  os.Exit(-1)
          //panic(err)
       }
    }
    command := os.Args[1]
    switch command {
	  case "enable":
                  enable()
          case "store":
		  scanner := bufio.NewScanner(os.Stdin)
		  scanner.Scan()
		  input := scanner.Text()
		  store(input)
	  case "get":
		  get()
	  case "erase":
		  erase()
	  default:
		  panic( errors.New("Vault Token Helper does not support command " + command))
    }

}
