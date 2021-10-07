package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireFileEqualAfterTrimSpace(t *testing.T, fileExpected, fileActual string) {
	contentExpected, errExpected := os.ReadFile(fileExpected)
	require.NoError(t, errExpected)

	contentActual, errActual := os.ReadFile(fileActual)
	require.NoError(t, errActual)

	require.Equal(t, strings.TrimSpace(string(contentExpected)), strings.TrimSpace(string(contentActual)))
}

func TestHosts(t *testing.T) {
	etcHosts = "/tmp/etcHosts"

	base, err := os.ReadFile("./tests/etc_hosts_base")
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(etcHosts, base, 0666))

	require.NoError(t, addHost("containerID1", "networkID1", "10.155.0.1", "service1.network.docker.internal"))
	requireFileEqualAfterTrimSpace(t, "./tests/etc_hosts_expected_0", etcHosts)

	require.NoError(t, addHost("containerID2", "networkID1", "10.155.0.2", "service2.network.docker.internal"))
	requireFileEqualAfterTrimSpace(t, "./tests/etc_hosts_expected_1", etcHosts)

	require.NoError(t, removeHostByID("containerID1", "networkID1"))
	requireFileEqualAfterTrimSpace(t, "./tests/etc_hosts_expected_2", etcHosts)

	require.NoError(t, removeHostByID("containerID2", "networkID1"))
	requireFileEqualAfterTrimSpace(t, "./tests/etc_hosts_expected_3", etcHosts)
}
