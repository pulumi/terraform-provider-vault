// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"fmt"
	"net/url"
)

func testAccDatabaseSecretsMount_mssql(name, path, pluginName string, parsedURL *url.URL) string {
	password, _ := parsedURL.User.Password()

	config := `
  mssql {
    allowed_roles     = ["dev", "prod"]
    plugin_name       = "%s"
    name              = "%s"
    connection_url    = "%s"
	username          = "%s"
	password          = "%s"
    verify_connection = true
  }`

	result := fmt.Sprintf(`
resource "vault_database_secrets_mount" "db" {
  path = "%s"
%s
}

resource "vault_database_secret_backend_role" "test" {
  backend = vault_database_secrets_mount.db.path
  name    = "dev"
  db_name = vault_database_secrets_mount.db.mssql[0].name
  creation_statements = [
    "CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';",
    "CREATE USER [{{name}}] FOR LOGIN [{{name}}];",
    "GRANT SELECT ON SCHEMA::dbo TO [{{name}}];",
  ]
}
`, path, fmt.Sprintf(config, pluginName, name, parsedURL.String(), parsedURL.User.Username(), password))

	return result
}

func testAccDatabaseSecretsMount_mssql_dual(name, name2, path, pluginName string, parsedURL *url.URL, parsedURL2 *url.URL) string {
	password, _ := parsedURL.User.Password()
	password2, _ := parsedURL2.User.Password()

	config := `
  mssql {
    allowed_roles     = ["dev1"]
    plugin_name       = "%s"
    name              = "%s"
    connection_url    = "%s"
	username          = "%s"
	password          = "%s"
    verify_connection = true
  }

  mssql {
    allowed_roles     = ["dev2"]
    plugin_name       = "%s"
    name              = "%s"
    connection_url    = "%s"
	username          = "%s"
	password          = "%s"
    verify_connection = true
  }
`
	result := fmt.Sprintf(`
resource "vault_database_secrets_mount" "db" {
  path = "%s"
%s
}

resource "vault_database_secret_backend_role" "test" {
  backend = vault_database_secrets_mount.db.path
  name    = "dev1"
  db_name = vault_database_secrets_mount.db.mssql[0].name
  creation_statements = [
    "CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';",
    "CREATE USER [{{name}}] FOR LOGIN [{{name}}];",
    "GRANT SELECT ON SCHEMA::dbo TO [{{name}}];",
  ]
}

resource "vault_database_secret_backend_role" "test2" {
  backend = vault_database_secrets_mount.db.path
  name    = "dev2"
  db_name = vault_database_secrets_mount.db.mssql[1].name
  creation_statements = [
    "CREATE LOGIN [{{name}}] WITH PASSWORD = '{{password}}';",
    "CREATE USER [{{name}}] FOR LOGIN [{{name}}];",
    "GRANT SELECT ON SCHEMA::dbo TO [{{name}}];",
  ]
}
`, path, fmt.Sprintf(config, pluginName, name, parsedURL.String(), parsedURL.User.Username(), password, pluginName,
		name2, parsedURL2.String(), parsedURL2.User.Username(), password2))

	return result
}
