# Security Model Specifications Language (SMSpec)

SMSpec is a YAML based language that specifies parts of a security model and connections between them.

## Security Model

A security model is made of three parts -

1. **External entities** - These are entities that are not part of a security analysis. We are only interested in their interactions with in-scope entities. No  security analysis is performed on these.
1. **In-scope entities** - These are entities that we want to analyze for security. We want to understand each of their characteristics and its impact on security.
1. **Flows** - These are communication links between entities (external or internal). For flows too, we want to understand their characteristics and security impact.

A security model starts with a header section

```yaml
design-document: "AnInvalidPath.md"
title: Sample security model
addb: "~/addb"
adm: ["adm/simple-addb.adm"]

```

* `design-document` - A link to a design document / file that was used as the source of information to build the security model.
* `title` - A name for the security model.
* `addb` - All required [ADDB](ADDB.md) entities are sourced from the directory specified under this field.
* `adm` - A list of ADM files that capture attacks and defenses for the entire security model. Typically these are items that span more than one entity and flow.

**NOTE**: ADDB location must be a directory on the local filesystem.

### External Entities

These are listed under the `externals` section of the model. Each entry can be a *human* or *program* specified using the `type` field.

```yaml
- id: user
  type: human
  name: Regular User
  description: A person interacting with the system. Only has basic access-rights to the system.
  interface: user-browser
```

* `id` - A unique identified for this entry. They are used in other parts of the model to refer to this entry.
* `type` - The type of entry. Valid values are `human` and `program`. In case of human, till someone invents the ultimate brain-computer interface, an additional `interface` field is required (as shown in example above) that points to a program used for interacting with the rest of the model items.
* `name` - A name for this entry. This is used by `adsm` tool for various purposes.
* `description` - One/Two line description about this entry.

### In-scope Entities

These are listed under the `entities` section of the model. Each entity can be a `human`, `program`, `system` or `role`. `system` and `program` are treated as synonyms.

#### Human

```yaml
- id: user
  type: human
  name: Alex
  description: Alex is a regular user. He uses the system as part of performing his duties.
  interface: alex-browser
  recommendations:
    - Always check the URL you are connecting to.
    - If you are expected to provide personal/private information on a website, always make sure it uses HTTPS URLs.
  adm: [alex.adm]
```

* `id` - A unique identified for this entry. They are used in other parts of the model to refer to this entry.
* `type` - The type of entry. Valid values are `human` and `program`. In case of human, an additional `interface` field is required (as shown in example above) that points to a program used for interacting with the rest of the model items.
* `name` - A name for this entry. This is used by `adsm` tool for various purposes.
* `description` - One/Two line description about this entry.
* `adm` - A list of ADM files that capture attacks targeted towards this entity and defenses that this entity must implement to mitigate those attacks.
* `recommendations` - A list of freeform recommendations for this entry. They must be generic and easy to understand for a non-technical reader.

#### Program / System

```yaml
- id: app-server
  type: program
  name: Application Server
  description: A service that receives requests from web-browser clients.
  repo: http://github.com/unknown/app-server
  base: addb:generic.application
  languages: [addb:lang.go, addb:lang.dockerfile]
  dependencies: [addb:libs.go.protobuf, addb:libs.go.http, addb:containers.alpine]
  roles: [login, regular-user]
  recommendations:
    - Always use HSTS for internet facing services.
    - Use least privilege access for each client to minimize attack surface.
  adm: [app-server.adm]
```

* `id` - A unique identified for this entry. They are used in other parts of the model to refer to this entry.
* `type` - The type of entry. Valid values are `human` and `program`. In case of human, an additional `interface` field is required (as shown in example above) that points to a program used for interacting with the rest of the model items.
* `name` - A name for this entry. This is used by `adsm` tool for various purposes.
* `description` - One/Two line description about this entry.
* `repo` - Applicable only to `program` entities. This field captures the link to the repository containing the code for this program and all necessary files required for its operation.
* `base` - Reference to another `program` entity spec. which is used as the base for this program. Base typically represents code/framework that this program is based on. Bases typically define a program's external-facing characteristics.
* `dependencies` - A list of references to `program` entities. Each of these entities may represent a library or software component that this program uses internally to meet its requirements. Examples include libraries like, protobuf, file-io, HTTP/TLS libraries, YAML/JSON/XML libraries, etc.
* `roles` - When this program plays a specific role when interacting with another program, a list of roles are specified. Each role specification captures attacks and defenses associated with the access granted to that role.
* `recommendations` - A list of freeform security recommendations for this entry. They must be generic and easy to understand for a non-technical reader.
* `adm` - A list of ADM files that capture attacks targeted towards this entity and defenses that this entity must implement to mitigate those attacks.

#### Role

Roles are lightweight entity specifications that capture risk associated with access rights to a system. Roles are associated with programs and may represent roles played by human via their interface or access rights granted to another program that makes requests to an entity.

```yaml
- id: add-record
  type: role
  name: Rights to add a record
  description: Role that limits access to adding a record only.
  adm: ["adm/add-record.adm"]
```

* `id` - A unique identified for this role. They are used in other parts of the model to refer to this entry.
* `type` - Roles are identified by the `role` type.
* `name` - A name for this role. This is used by `adsm` tool for various purposes.
* `description` - One/Two line description about this role.
* `adm` - A list of ADM files that capture attacks targeted towards this role and defenses that must be implemented to mitigate them.

### Flow

```yaml
- id: process-requests
  name: Process requests
  description: System processing user requests
  sender: frontend
  receiver: backend
  protocol: [addb:flow.https] 
  adm: ["adm/add-request.adm", "adm/update-request.adm", "adm/delete-request.adm", "adm/read-request.adm"]
```

* `id` - A unique identified for this flow. They are used in other parts of the model to refer to this entry.
* `name` - A name for this role. This is used by `adsm` tool for various purposes.
* `description` - One/Two line description about this role.
* `sender` - ID of the entity sending the request.
* `receiver` - ID of the entity receiving the request. It is also the one that will send the response back to the `sender`.
* `protocol` - The communication mechanism used for this flow. This is specified as a protocol stack with each entry pointing to the ID of a specific protocol.
* `adm` - A list of ADM files that capture attacks targeted towards this flow and defenses that must be implemented to mitigate them.

*NOTE: Flows don't have a `type` field.*

## YAML schema

This repository contains a schema specification - `schemas/model-schema.json` that can be used when building a security model. If you add `yaml-language-server: $schema= [PATH_TO_MODEL_SCHEMA_JSON]` as the first line of the YAML file, a text-editor / IDE that supports YAML Language Server will use it to validate your model's structure.

## ADDB

References to ADDB entities can be included in a security model using `addb:<entity-id>` in any field that references an ID (like `language`, `base`, `dependencies`, `protocol`, etc.).
