package tool

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)


func NewReader() *bufio.Reader {
	reader := bufio.NewReader(os.Stdin)
	return reader
}


func ReadStdin(reader *bufio.Reader, describe string)(string, error) {
	fmt.Print(describe)
	readString, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	readString = strings.Replace(readString, "\n", "", 1)
	return readString, nil
}


func Input(describe string, defaultString string) (string, error) {
	reader := NewReader()
	readStdin, err := ReadStdin(reader, describe)
	if err != nil {
		fmt.Println(err)
		return  "", err
	}
	readStdin = strings.Trim(readStdin, "\r")
	if readStdin == "" {
		readStdin = defaultString
	}
	fmt.Println(readStdin)
	return readStdin, nil
}


func Confirm(describe string, retry int) (bool, error)  {
	describe += " please enter (y/n):"

	i := 1
	for i <= retry {
		isConfirm, err := confirm(describe)
		if err != nil{
			fmt.Println(err)
			i++
			continue
		}
		return isConfirm, nil
	}
	return false, nil
}


func confirm(describe string) (bool, error){
	reader := NewReader()
	isConfirm, err := ReadStdin(reader, describe)
	if err != nil{
		return false, err
	}
	isConfirm = strings.ToLower(strings.Trim(isConfirm, "\r"))
	if  isConfirm == "n" || isConfirm == "no"{
		return false, nil
	}
	if  isConfirm == "y" || isConfirm == "yes"{
		return true, nil
	}
	return false, errors.New("please enter (y/n)")
}

func ReadLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}