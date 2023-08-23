// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/internal/provider"
	"github.com/hashicorp/terraform-provider-vault/testutil"
)

type testMountStore struct {
	uuid string
	path string
}

func TestConsulSecretBackend(t *testing.T) {
	t.Parallel()
	path := acctest.RandomWithPrefix("tf-test-consul")
	resourceType := "vault_consul_secret_backend"
	resourceName := resourceType + ".test"
	token := "026a0c16-87cd-4c2d-b3f3-fb539f592b7e"

	resource.Test(t, resource.TestCase{
		Providers:    testProviders,
		PreCheck:     func() { testutil.TestAccPreCheck(t) },
		CheckDestroy: testCheckMountDestroyed(resourceType, consts.MountTypeConsul, consts.FieldPath),
		Steps: []resource.TestStep{
			{
				Config: testConsulSecretBackend_initialConfig(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "address", "127.0.0.1:8500"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "http"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldLocal, "false"),
					resource.TestCheckNoResourceAttr(resourceName, "ca_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "client_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "client_key"),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil,
				"token", "bootstrap", "ca_cert", "client_cert", "client_key", "disable_remount"),
			{
				Config: testConsulSecretBackend_initialConfigLocal(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "address", "127.0.0.1:8500"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "http"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldLocal, "true"),
					resource.TestCheckNoResourceAttr(resourceName, "ca_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "client_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "client_key"),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil,
				"token", "bootstrap", "ca_cert", "client_cert", "client_key", "disable_remount"),
			{
				Config: testConsulSecretBackend_updateConfig(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "address", "consul.domain.tld:8501"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "https"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldLocal, "false"),
					resource.TestCheckNoResourceAttr(resourceName, "ca_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "client_cert"),
					resource.TestCheckNoResourceAttr(resourceName, "client_key"),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil,
				"token", "bootstrap", "ca_cert", "client_cert", "client_key", "disable_remount"),
			{
				Config: testConsulSecretBackend_updateConfig_addCerts(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "address", "consul.domain.tld:8501"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "https"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldLocal, "false"),
					resource.TestCheckResourceAttr(resourceName, "ca_cert", "FAKE-CERT-MATERIAL"),
					resource.TestCheckResourceAttr(resourceName, "client_cert", "FAKE-CLIENT-CERT-MATERIAL"),
					resource.TestCheckResourceAttr(resourceName, "client_key", "FAKE-CLIENT-CERT-KEY-MATERIAL"),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil,
				"token", "bootstrap", "ca_cert", "client_cert", "client_key", "disable_remount"),
			{
				Config: testConsulSecretBackend_updateConfig_updateCerts(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, consts.FieldPath, path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "0"),
					resource.TestCheckResourceAttr(resourceName, "address", "consul.domain.tld:8501"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "https"),
					resource.TestCheckResourceAttr(resourceName, consts.FieldLocal, "false"),
					resource.TestCheckResourceAttr(resourceName, "ca_cert", "FAKE-CERT-MATERIAL"),
					resource.TestCheckResourceAttr(resourceName, "client_cert", "UPDATED-FAKE-CLIENT-CERT-MATERIAL"),
					resource.TestCheckResourceAttr(resourceName, "client_key", "UPDATED-FAKE-CLIENT-CERT-KEY-MATERIAL"),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil,
				"token", "bootstrap", "ca_cert", "client_cert", "client_key", "disable_remount"),
		},
	})
}

func TestConsulSecretBackend_remount(t *testing.T) {
	t.Parallel()
	path := acctest.RandomWithPrefix("tf-test-consul")
	updatedPath := acctest.RandomWithPrefix("tf-test-consul-updated")
	token := "026a0c16-87cd-4c2d-b3f3-fb539f592b7e"

	resourceType := "vault_consul_secret_backend"
	resourceName := resourceType + ".test"

	store := &testMountStore{}

	resource.Test(t, resource.TestCase{
		Providers:    testProviders,
		PreCheck:     func() { testutil.TestAccPreCheck(t) },
		CheckDestroy: testCheckMountDestroyed(resourceType, consts.MountTypeConsul, consts.FieldPath),
		Steps: []resource.TestStep{
			{
				Config: testConsulSecretBackend_initialConfig(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "path", path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "address", "127.0.0.1:8500"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "http"),
					testCaptureMountUUID(path, store),
				),
			},
			{
				Config: testConsulSecretBackend_initialConfig(updatedPath, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "path", updatedPath),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "address", "127.0.0.1:8500"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "http"),
					testMountCompareUUIDs(updatedPath, store, true),
					testCaptureMountUUID(updatedPath, store),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil, "token", "bootstrap", "disable_remount"),
			{
				Config: testConsulSecretBackend_disableRemount(path, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "path", path),
					resource.TestCheckResourceAttr(resourceName, consts.FieldDescription, "test description"),
					resource.TestCheckResourceAttr(resourceName, "default_lease_ttl_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "max_lease_ttl_seconds", "86400"),
					resource.TestCheckResourceAttr(resourceName, "address", "127.0.0.1:8500"),
					resource.TestCheckResourceAttr(resourceName, "token", token),
					resource.TestCheckResourceAttr(resourceName, "scheme", "http"),
					resource.TestCheckResourceAttr(resourceName, "disable_remount", "true"),
					testMountCompareUUIDs(path, store, false),
				),
			},
		},
	})
}

func testCaptureMountUUID(path string, store *testMountStore) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		mount, err := testGetMount(path)
		if err != nil {
			return err
		}

		store.path = path
		store.uuid = mount.UUID

		return nil
	}
}

func testMountCompareUUIDs(path string, store *testMountStore, equal bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		mount, err := testGetMount(path)
		if err != nil {
			return err
		}

		if store.uuid == mount.UUID {
			if !equal {
				return fmt.Errorf("expected different uuids after mount creation; "+
					"both uuids equal to %s", store.uuid)
			}
		} else {
			if equal {
				return fmt.Errorf("expected same uuid after remount; "+
					"got id1=%s, id2=%s", store.uuid, mount.UUID)
			}
		}

		return nil
	}
}

func testGetMount(path string) (*api.MountOutput, error) {
	client, err := provider.GetClient("", testProvider.Meta())

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	mount, ok := mounts[strings.Trim(path, "/")+"/"]
	if !ok {
		return nil, fmt.Errorf("given mount %s not found", path)
	}

	return mount, nil
}

func testConsulSecretBackend_initialConfig(path, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path = "%s"
  description = "test description"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds = 86400
  address = "127.0.0.1:8500"
  token = "%s"
}`, path, token)
}

func testConsulSecretBackend_disableRemount(path, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path = "%s"
  description = "test description"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds = 86400
  address = "127.0.0.1:8500"
  token = "%s"
  disable_remount = true
}`, path, token)
}

func testConsulSecretBackend_initialConfigLocal(path, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path = "%s"
  description = "test description"
  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds = 86400
  address = "127.0.0.1:8500"
  token = "%s"
  local = true
}`, path, token)
}

func testConsulSecretBackend_bootstrapConfig(path, addr, token string, bootstrap bool) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path = "%s"
  description = "test description"
  address = "%s"
  token = "%s"
  bootstrap = %t
  disable_remount = false
}
`, path, addr, token, bootstrap)
}

func testConsulSecretBackend_updateConfig(path, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path        = "%s"
  description = "test description"
  address     = "consul.domain.tld:8501"
  token       = "%s"
  scheme      = "https"
}
`, path, token)
}

func testConsulSecretBackend_bootstrapAddRole(path, addr string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path            = "%s"
  description     = "test description"
  address         = "%s"
  bootstrap       = true
  disable_remount = false
}

resource "vault_consul_secret_backend_role" "test" {
  backend         = vault_consul_secret_backend.test.path
  name            = "management"
  consul_policies = ["global-management"]
}
`, path, addr)
}

func testConsulSecretBackend_bootstrapAddRoleMulti(path, addr string) string {
	return fmt.Sprintf(`
%s

resource "vault_consul_secret_backend" "test-2" {
  path            = "%s-2"
  description     = "test description"
  address         = "%s"
  bootstrap       = true
  disable_remount = false
}
`, testConsulSecretBackend_bootstrapAddRole(path, addr), path, addr)
}

func testConsulSecretBackend_updateConfig_addCerts(path, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path = "%s"
  description = "test description"
  address = "consul.domain.tld:8501"
  token = "%s"
  scheme = "https"
  ca_cert = "FAKE-CERT-MATERIAL"
  client_cert = "FAKE-CLIENT-CERT-MATERIAL"
  client_key = "FAKE-CLIENT-CERT-KEY-MATERIAL"
}`, path, token)
}

func testConsulSecretBackend_updateConfig_updateCerts(path, token string) string {
	return fmt.Sprintf(`
resource "vault_consul_secret_backend" "test" {
  path = "%s"
  description = "test description"
  address = "consul.domain.tld:8501"
  token = "%s"
  scheme = "https"
  ca_cert = "FAKE-CERT-MATERIAL"
  client_cert = "UPDATED-FAKE-CLIENT-CERT-MATERIAL"
  client_key = "UPDATED-FAKE-CLIENT-CERT-KEY-MATERIAL"
}`, path, token)
}
