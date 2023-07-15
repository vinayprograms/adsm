# ADDB - Attack-Defense Database

ADDB is a repository of common entity specifications that can then be reused in multiple security models. It is made up of a set of entity files and associated ADM files. Entity and ADM files can be grouped into directories based on parameters like company name, technology used, programming language, etc.

Each entity specification contains information about that entity, similar to those specified in security models, including the list of ADM files for that entity. A file can contain more than one entity. These are separated via `---` and `...` YAML delimiters.

## Structure

```yaml
# yaml-language-server: $schema=../component-schema.json
---
id: lang.html
type: program
name: HTML (HyperText Markup Language)
description: ADM specifications for HTML
repo: nil
languages: []
adm: [html.adm]
...

---
id: lang.javascript
name: Javascript
description: ADM specifications for javascript programming language
type: program
repo: nil
languages: []
adm: [javascript.adm]
...

---
id: lib.py.torch@1.6.0
type: program
name: Torch (python) v1.6.0
description: Tensors and Dynamic neural networks in Python with strong GPU acceleration. This specification is for v1.6.0
repo: https://github.com/pytorch/pytorch/tree/v1.6.0
languages: [python, cpp, cuda, c, objc, cmake]
adm: [torch@1.6.0.adm]
...

```

An ADDB entry must have the following mandatory fields -

* `id` - A unique identified for this entry. IDs can be specified with a N-level namespace pattern - `<root>.<sub-group>.<id>`. ID must always be unique within ADDB. The ID is independent of the directory structure used in ADDB. If you want to include version information (typically required for libraries & frameworks), you can append a `@<VERSION>`.
* `name` - A name for this entry. This is used by `adsm` tool for various purposes.
* `description` - One/Two line description about this ADDB entry.
* `type` - The type of entry. Valid values are - `human`, `program`, `role` and `flow`.
* `adm` - A list of ADM files that capture attacks targeted towards this entity and defenses that this entity must implement to mitigate those attacks.

## Fields for each entity type

In addition to the mandatory ones, each type of entity can have the following additional fields.

### Human

* `base` - Reference to another `human` entity spec. which is used as the base for this entity. Base typically represents generic human specification that this spec. is based on.
* `interface` - Reference to the program used by the human to interact with the system. Depending on the nature of the system, this can be a web-browser or some client application.
* `recommendations` - A list of freeform security recommendations for this human. They must be generic and easy to understand for a non-technical reader.

### Program

* `repo` - A link to the repository containing the code for this program. Depending on the specification, it can point to the root of the repository or a specific branch, commit or tag.
* `base` - Reference to another `program` entity spec. which is used as the base for this program. Base typically represents the generic pattern that this program is based on. Examples include `web server`, `headless server`, `iOS app`, `Android app`, etc. Bases define the inherent characteristics a program should have to operate in its target environment.
* `languages` - A list of references to `program` entities. Each of these may represent a programming language used in the program. Apart from a primary language used to write the program's logic, it may use other languages to meet its operational and deployment needs. Examples include SQL, shell scripts, etc.
* `dependencies` - A list of references to `program` entities. Each of these may represent a library or software component that this program uses internally to meet its requirements. Examples include libraries like, protobuf, file-io, HTTP/TLS libraries, YAML/JSON/XML libraries, etc.
* `roles` - If this program plays a specific role when interacting with another program, a list of roles are specified. Each role specification captures attacks and defenses associated with the access granted to that role. The most common use of roles is with user-interface programs. After a human successfully logs in, their user-interface program is assigned one or more roles representing all the access rights assigned to that human. Roles can also be used between two programs if strict access control is implemented for interactions within the system.
* `recommendations` - A list of freeform security recommendations for this entry. They must be generic and easy to understand for a non-technical reader.

### Role

A list of security `recommendations` can be added to a role in ADDB. It is a good idea to add *role* entities to ADDB when you use common role specifications across multiple products / designs.

## Flow

* `protocol` - The communication mechanism used for this flow. This is specified as a protocol stack with each entry pointing to the ID of a specific protocol.
* `recommendations` - A list of freeform security recommendations for this entry. They must be generic and easy to understand for a non-technical reader.
