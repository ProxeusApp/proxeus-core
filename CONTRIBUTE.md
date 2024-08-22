---
title: Developer docs for contributors
---

## Hello, world!

:+1::tada: First, thanks for getting involved with Proxeus! :tada::+1:

https://github.com/ProxeusApp 

Structure of **[proxeus-core](https://github.com/ProxeusApp/proxeus-core)**:

- **demo**: database dumps with preconfigured content for a fresh or demo server
- **deploy**: initialization scripts for supported cloud hosting providers
- **docs**: documentation in Markdown formatting
- **doc**: auto-generated license documentation
- **externalnode**: library for plugin nodes
- **main**: the Go backend and web application main source tree
- **service**: definition of external API services
- **storage**: database and file handling
- **sys**: key utilities, libraries, models, internal services
- **test**: unit test cases for automated backend & end-to-end testing

Other key repositories include:

- **[document-service](https://github.com/ProxeusApp/document-service)**: handling documents, generating PDFs
- **[node-go](https://github.com/ProxeusApp/node-go)**: Golang library for nodes (Proxeus extensions)
- **[proxeus-contract](https://github.com/ProxeusApp/proxeus-contract)**: Solidity code for blockchain support

## How to contribute

You can start by checking our [core issues board](https://github.com/ProxeusApp/proxeus-core/issues). A few of them are labelled `good first issue` or `hacktoberfest` or `bounty`!

If you would like to work on something, you can:

1. Comment on the issue to let us know you'd like to work on, so that we can assist you and to make sure no one has started looking into it yet.
2. If good to go, fork the main repository.
3. Clone the forked repository to your machine.
4. Create a feature branch (e.g. `50-update-readme`, where `50` is the number of the related issue).
5. Commit your changes to the feature branch.
6. Add comments to the issue describing the changes.
7. Push the feature branch to your forked repository.
8. Create a Pull Request against the original repository.
   - add a short description of the changes included in the PR
9. Address review comments if requested by our demanding reviewers 😜.

If you have an idea for improvement, and it doesn't have a corresponding issue yet, simply submit a new one.

> Join our [GitHub Discussions](https://github.com/orgs/ProxeusApp/discussions) to discuss existing issues and to ask for help.

