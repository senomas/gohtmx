package mariadb_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func startMariaDB(t *testing.T) {
	out, err := os.Create("mariadb.log")
	assert.NoError(t, err)
	defer out.Close()

	fmt.Println("docker rm -f account_store_test")
	cmd := exec.Command("docker", "rm", "-f", "account_store_test")
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	if os.Getenv("DOCKER_PULL") != "" {
		fmt.Println("docker pull mariadb:lts")
		cmd = exec.Command("docker", "pull", "mariadb:lts")
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("docker run --name account_store_test mariadb:lts")
	cmd = exec.Command(
		"docker",
		"run",
		"--name", "account_store_test",
		"-p", "13306:3306",
		"-e", "MARIADB_ROOT_PASSWORD=dodol123",
		"-e", "MARIADB_DATABASE=test",
		"mariadb:lts",
	)
	cmd.Stdout = out
	cmd.Stderr = out
	go func() {
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("docker exec account_store_test mysqladmin ping")
	for i := 0; i < 60; i++ {
		cmd = exec.Command(
			"docker", "exec", "account_store_test",
			"mysqladmin", "ping", "-pdodol123", "-uroot")
		out, err := cmd.CombinedOutput()
		str := string(out)
		if err != nil {
			fmt.Println(err)
		} else {
			if str == "mysqld is alive\n" {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func stopMariaDB(t *testing.T) {
	cmd := exec.Command("docker", "rm", "-f", "account_store_test")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}
