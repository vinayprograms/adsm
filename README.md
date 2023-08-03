# ADSM - [ADM](https://github.com/vinayprograms/adm) based Security Modeling

This is a CLI tool that processes `.smspec` files written using [smspec YAML schema](SMSPEC.md).

This tool can -

1. Associate attack/defense specification written in ADM language with each entity and flow in a security model.
1. Generate consolidated statistics about each entity and flow.
1. Generate ADM and security model diagrams.
1. Generate markdown report listing all security recommendations and unmitigated security risks.

## Building from source

* Clone this repository and `adm` repository into a single parent directory. ADM repository is used by ADSM to handle attack-defense files.
* Run `make build`. This will create a `adsm` binary in `bin` directory.

```shell
> mkdir adm-adsm && cd adm-adsm
> git clone [adm-git-url]
> git clone [adsm-git-url]
> cd adsm && make build
```

If you want the tool to be available from anywhere in the system, it is a good practice to copy the binary to a `bin` folder in home directory - `mkdir ~/.local/bin` and append this path to `$PATH` variable in your `.bashrc`/`.zshrc` file.

## Usage

***NOTE**: You can use examples in `test/examples/` to experiment with commands shown in this section.*

If no path is specified, `adsm` will print a quick help showing all the subcommands and associated flags.

### `stat` sub-command

1. Without a flag - `adsm stat [path to .smspec file]`, will list all entities and flows along with a single line summary about each associated ADM file.
1. Following flags can be used to filter output -
    * `-x` - Only list external entities. For example output of `./bin/adsm stat -x test/examples/simple_addb.smspec` will be

      ```text
      MODEL: Sample security model
        External Entity: Regular User
                         A person interacting with the system. Only has basic access-rights to the system.
        External Entity: Administrator
                         Administrator of the system. Controls system's configuration and operations.
      ```

    * `-e` - Only show statistics for entities in model's scope. For example output of `./bin/adsm stat -e test/examples/simple_addb.smspec` will be

      ```text
      MODEL: Sample security model
        Entity: User's web-browser
                Stock web-browser running on user's device.
                ADM: test/examples/adm/modify-record.adm, ATTACKS:1, DEFENSES:1
                ADM: test/examples/adm/read-record.adm, ATTACKS:1, DEFENSES:1
                ADM: test/examples/adm/add-record.adm, ATTACKS:1, POLICIES:1
        Entity: Administrator's web-browser
                Stock web-browser running on administrator's device.
                ADM: test/examples/adm/modify-record.adm, ATTACKS:1, DEFENSES:1
                ADM: test/examples/adm/read-record.adm, ATTACKS:1, DEFENSES:1
                ADM: test/examples/adm/delete-record.adm, ATTACKS:1, POLICIES:1
                ADM: test/examples/adm/add-record.adm, ATTACKS:1, POLICIES:1
        Entity: Web UI
                A Web UI with which both administrators and regular users interact.
                ADM: test/examples/adm/frontend.adm, ASSUMPTIONS:1, ATTACKS:1
        Entity: Business logic
                Server that processes all requests
                ADM: test/examples/adm/backend.adm, ATTACKS:1, DEFENSES:1
        Entity: Database
                Database used by business-logic to persist important data
                ADM: test/examples/adm/db.adm, ATTACKS:1, POLICIES:1
      ```

    * `-r` - Only show statistics for roles. For example output of `./bin/adsm stat -r test/examples/simple_addb.smspec` will be

      ```text
      MODEL: Sample security model
        Role: Rights to read a record
              Role that limits access to reading a record only.
              ADM: test/examples/adm/read-record.adm, ATTACKS:1, DEFENSES:1
        Role: Rights to delete a record
              Role that limits access to deleting a record only.
              ADM: test/examples/adm/delete-record.adm, ATTACKS:1, POLICIES:1
        Role: Rights to add a record
              Role that limits access to adding a record only.
              ADM: test/examples/adm/add-record.adm, ATTACKS:1, POLICIES:1
        Role: Rights to modify a record
              Role that limits access to modifying a record only.
              ADM: test/examples/adm/modify-record.adm, ATTACKS:1, DEFENSES:1
      ```

    * `-f` - Only show flows. For example output of `./bin/adsm stat -f test/examples/simple_addb.smspec` will be

      ```text
      MODEL: Sample security model
        Flow: User Interaction
              User's interaction with the frontend
              ADM: /Users/vinay/addb/sm/flows/adm/tls.adm, ATTACKS:1, DEFENSES:2
        Flow: Administrator Interaction
              Administrator's interaction with the frontend
              ADM: /Users/vinay/addb/sm/flows/adm/tls.adm, ATTACKS:1, DEFENSES:2
        Flow: Process requests
              System processes user request
              ADM: test/examples/adm/add-request.adm, ATTACKS:1, POLICIES:1
              ADM: test/examples/adm/update-request.adm, ATTACKS:1, DEFENSES:2
              ADM: test/examples/adm/delete-request.adm, ATTACKS:1, DEFENSES:1
              ADM: test/examples/adm/read-request.adm, ATTACKS:1, DEFENSES:1
              ADM: /Users/vinay/addb/sm/flows/adm/tls.adm, ATTACKS:1, DEFENSES:2
        Flow: Insert or update data
              Update or Insert user data. User may imply regular user or administrator.
              ADM: test/examples/adm/add-db.adm, ATTACKS:1, POLICIES:1
              ADM: test/examples/adm/update-db.adm, ATTACKS:1, DEFENSES:1
      ```

### `diag` sub-command

This subcommand generates two diagrams

* Security model - generates a single graphviz file showing entities and flows
* ADM - Generates a single graphviz file containing all the ADM graphs

For example, `adsm diag test/examples/simple_addb.smspec` generates both the diagrams and places them in the current directory. The target directory can be specified using `-d` flag - `adsm diag -d ~/reports test/examples/simple_addb.smspec`.
If you need only one of the diagrams, pass `-sm` or `-adm` flag to the command. For example `adsm diag -sm -d ~/reports test/examples/simple_addb.smspec` will only output the security model as a graphviz file.

### `report` sub-command

This subcommand generates a markdown file containing

1. Security Model diagram
1. Consolidated list of recommendations for specific entities and flows
1. A list of un-mitigated risks for specific entities and flows

The report (and associated diagram) is written to a `/report` folder in the current directory. You can change the location using `-d` flag. For example `adsm report -d ~/smreports test/examples/simple_addb.smspec` will create a `report` subdirectory under `~/smreports`.

## ADDB

Security model entities / flows can be reused by adding them to a *Attack-Defense Database*. This is a git repository containing entity specifications along with its ADM files. See [ADDB](ADDB.md) for more details.
