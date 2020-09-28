---
layout: "flexibleengine"
page_title: "flexibleengine: flexibleengine_cce_addon_v3"
sidebar_current: "docs-flexibleengine-resource-cce-addon_v3"
description: |-
  Add an addon to a container cluster. 
---

# flexibleengine_cce_addon_v3

Provides a addon resource (CCE).


## Example Usage
```hcl
variable "cluster_id" { }

resource "flexibleengine_cce_addon_v3" "addon_test" {
    cluster_id    = var.cluster_id
    template_name = "metrics-server"
    version       = "1.0.0"
}
``` 

## Argument Reference
The following arguments are supported:
* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.
* `template_name` - (Required) Name of the addon template. Changing this parameter will create a new resource.
* `version` - (Required) Version of the addon. Changing this parameter will create a new resource.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

 * `id` -  ID of the addon instance.
 * `status` - Addon status information.
 * `description` - Description of addon instance.
