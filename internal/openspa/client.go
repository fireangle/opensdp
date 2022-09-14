package openspa

import (
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	Cmd    string
	OSPA   string
	Server net.IP
	Port   uint16
}

type Request struct {
	Protocol  string
	StartPort uint16
	EndPort   uint16
}

func (c *Client) Send(req Request, continuous bool, clientIp string) error {

	sPort := strconv.Itoa(int(req.StartPort))
	ePort := strconv.Itoa(int(req.EndPort))
	serverPort := strconv.Itoa(int(c.Port))

	cmdStr := []string{c.Cmd, "request", c.OSPA, "--protocol", req.Protocol, "-p", sPort, "--end-port", ePort,
		"--server-ip", c.Server.String(), "--server-port", serverPort, "--client-ip", clientIp}

	if continuous {
		cmdStr = append(cmdStr, "-a")
	}

	log.WithField("command", strings.Join(cmdStr, " ")).Debug("OpenSPA command")

	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)

	// Get stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		return err
	}

	defer cmd.Wait()
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return nil
}
