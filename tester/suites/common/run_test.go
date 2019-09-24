package common

import (
	"testing"
)

func TestRunSuite(t *testing.T) {
	SetUp()
	defer TearDown()
	runCase(t, TestCheckIPList)
	runCase(t, TestTCPCheckIPList)
	runCase(t, TestURLWrite)
	runCase(t, TestRoundRobin)
	runCase(t, TestAccessControl)
	runCase(t, TestFlowControl)
	runCase(t, TestTCPFlowControl)
	runCase(t, TestAppAccess)
	runCase(t, TestFilter)
}

func runCase(t *testing.T, testCase func(*testing.T)) {
	Before()
	defer After()
	testCase(t)
}