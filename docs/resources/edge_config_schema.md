---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vercel_edge_config_schema Resource - terraform-provider-vercel"
subcategory: ""
description: |-
  An Edge Config Schema provides an existing Edge Config with a JSON schema. Use schema protection to prevent unexpected updates that may cause bugs or downtime.
---

# vercel_edge_config_schema (Resource)

An Edge Config Schema provides an existing Edge Config with a JSON schema. Use schema protection to prevent unexpected updates that may cause bugs or downtime.

## Example Usage

```terraform
resource "vercel_edge_config" "example" {
  name = "example"
}

resource "vercel_edge_config_schema" "example" {
  id         = vercel_edge_config.example.id
  definition = <<EOF
{
  "title": "Greeting",
  "type": "object",
  "properties": {
    "greeting": {
      "description": "A friendly greeting",
      "type": "string"
    }
  }
}
EOF
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `definition` (String) A JSON schema that will be used to validate data in the Edge Config.
- `id` (String) The ID of the Edge Config that the schema should apply to.

### Optional

- `team_id` (String) The ID of the team the Edge Config should exist under. Required when configuring a team resource if a default team has not been set in the provider.

## Import

Import is supported using the following syntax:

```shell
# If importing into a personal account, or with a team configured on
# the provider, simply use the edge config id.
# - edge_config_id can be found by navigating to the Edge Config in the Vercel UI. It should begin with `ecfg_`.
terraform import vercel_edge_config.example ecfg_xxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Alternatively, you can import via the team_id and edge_config_id.
# - team_id can be found in the team `settings` tab in the Vercel UI.
# - edge_config_id can be found by navigating to the Edge Config in the Vercel UI. It should begin with `ecfg_`.
terraform import vercel_edge_config.example team_xxxxxxxxxxxxxxxxxxxxxxxx/ecfg_xxxxxxxxxxxxxxxxxxxxxxxxxxxx
```