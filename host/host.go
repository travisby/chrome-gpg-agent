package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Message struct {
	Method    string  `json:"method"`
	Message   string  `json:"message"`
	Recipient *string `json:"recipient"`
}

type Output struct {
	Message string `json:"message"`
}

var ERR_WRONG_BYTES = errors.New("Incorrect number of bytes read.  Expected 4")

func main() {
	msg, err := readMessage(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var message *string
	switch msg.Method {
	case "sign":
		message, err = sign(msg.Message)
	case "verify":
		message, err = verify(msg.Message)
	case "encrypt":
		// TODO pointer check
		message, err = encrypt(msg.Message, *msg.Recipient)
	case "decrypt":
		message, err = decrypt(msg.Message)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := writeMessage(Output{*message}, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func sign(str string) (*string, error) {
	cmd := exec.Command("gpg", "--clearsign")
	cmd.Stdin = bytes.NewReader([]byte(str))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	temp := string(output)
	return &temp, nil
}

func verify(str string) (*string, error) {
	// verify is unique.  the interesting bits actually happen on stderr
	cmd := exec.Command("gpg", "--verify")
	cmd.Stdin = bytes.NewReader([]byte(str))
	cmd.Stdout = nil

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	temp := string(output)
	return &temp, nil
}

func encrypt(str, recipient string) (*string, error) {
	cmd := exec.Command("gpg", "--encrypt", "-r", recipient, "--armor")
	cmd.Stdin = bytes.NewReader([]byte(str))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	temp := string(output)
	return &temp, nil
}

func decrypt(str string) (*string, error) {
	cmd := exec.Command("gpg", "--decrypt")
	cmd.Stdin = bytes.NewReader([]byte(str))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	temp := string(output)
	return &temp, nil
}

func writeMessage(message interface{}, w io.Writer) error {
	sendOnWire, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, int32(len(sendOnWire))); err != nil {
		return err
	}

	if n, err := w.Write(sendOnWire); err != nil {
		return err
	} else if n != len(sendOnWire) {
		return fmt.Errorf("Expected to write %d bytes but actually wrote %d", len(sendOnWire), n)
	}

	return nil
}

func readMessage(rd io.Reader) (*Message, error) {
	var msgSize int32
	if err := binary.Read(rd, binary.LittleEndian, &msgSize); err != nil {
		return nil, err
	}

	var msg Message
	if err := json.NewDecoder(rd).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
