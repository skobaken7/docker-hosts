package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/juju/fslock"
)

const markerFormat = "# added by docker-hosts; container=%s, network=%s"

var etcHosts = "/etc/hosts"

func addHost(containerID, networkID, ip, fqdn string) error {
	log.Printf("add %s <- %s\n", ip, fqdn)

	marker := fmt.Sprintf(markerFormat, containerID, networkID)

	lock := fslock.New(etcHosts)
	if err := lock.Lock(); err != nil {
		return errors.Wrapf(err, "failed to lock %s", etcHosts)
	}
	defer lock.Unlock()

	f, err := os.OpenFile(etcHosts, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", etcHosts)
	}

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return errors.Wrapf(err, "failed to read line of %s", etcHosts)
		}

		if strings.Contains(line, marker) {
			return nil
		}

		if err == io.EOF {
			break
		}
	}

	f.WriteString("\n")
	f.WriteString(fmt.Sprintf("%s\t%s\t%s", ip, fqdn, marker))
	return nil
}

func removeHostByID(containerID, networkID string) error {
	log.Printf("remove hosts by containerID=%s,networkID=%s\n", containerID, networkID)

	marker := []byte(fmt.Sprintf(markerFormat, containerID, networkID))

	lock := fslock.New(etcHosts)
	if err := lock.Lock(); err != nil {
		return errors.Wrapf(err, "failed to lock %s", etcHosts)
	}
	defer lock.Unlock()

	contents, err := os.ReadFile(etcHosts)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", etcHosts)
	}

	updated := make([]byte, 0, len(contents))
	reader := bufio.NewReader(bytes.NewReader(contents))
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return errors.Wrapf(err, "failed to read line of %s", etcHosts)
		}

		if !bytes.Contains(line, marker) {
			updated = append(updated, line...)
		}

		if err == io.EOF {
			break
		}
	}

	if err := os.WriteFile(etcHosts, updated, 0644); err != nil {
		return errors.Wrapf(err, "failed to read %s", etcHosts)
	}

	return nil
}
