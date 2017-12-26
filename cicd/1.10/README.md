# Continuous Integration and Delivery with GitLab and Jenkins

This demo presents a continuous integration and delivery (CI/CD) workflow using GitLab and Jenkins on DC/OS. In this demo we’ll be creating a simple website, testing it to make sure all the links work, and then deploying it in production if it passes.

- Estimated time for completion: 30min
- Target audience: Anyone interested in simplifying the deployment of a CI/CD pipeline.

**Table of Contents**:

- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Install](#install)
- [Use the demo](#use)

## Architecture

## Prerequisites

- A running [DC/OS 1.10](https://dcos.io/releases/) or higher cluster with at least 3 private agents and 1 public agent. Each agent should have 2 CPUs and 5 GB of RAM available. The [DC/OS CLI](https://dcos.io/docs/1.10/usage/cli/install/) also needs to be installed.
- The [dcos/demo](https://github.com/dcos/demos/) Git repo must be available locally, use: `git clone https://github.com/dcos/demos/` if you haven't done so yet.
- [SSH](https://dcos.io/docs/1.10/administration/access-node/sshcluster/) cluster access must be set up.

The DC/OS services used in the demo are as follows:

- GitLab
- Jenkins
- Marathon-lb
- Universal Container Runtime (UCR)

## Install

Any configuration needed accept images from an insecure registry

Install Marathon-LB

Install GitLab

Install Jenkins

Configure job in Jenkins

- Support custom Git repository for website HTML files
- Build website Docker image
- Upload to GitLab Docker registry
- UCR to launch the Docker image for testing with http://wummel.github.io/linkchecker/
- Once passed, launch Docker image with UCR into production

### 

Should you have any questions or suggestions concerning the demo, please raise an [issue](https://jira.mesosphere.com/) in Jira or let us know via the [users@dcos.io](mailto:users@dcos.io) mailing list.

