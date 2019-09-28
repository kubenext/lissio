# Frequently Asked Questions

## Q: How is Lissio different from the _"official"_ Kubernetes Dashboard?

Kubernetes Dashboard is described as a **_"general purpose_**, web-based UI for Kubernetes clusters.", whereas **Lissio** is designed to be a _"tool for **developers** to understand how applications run on a Kubernetes cluster."_ 

Lissio provides more detail, is more extensible, uses newer technology, and is under more active development.

More specifically:
- Lissio does not run in a cluster (by default). Instead, it runs locally on your workstation and uses your kubeconfig files
- Lissio has a **resource viewer** which links related objects to better describe their relationship within the cluster
- Lissio supports Custom Resource Definitions (CRDs)
- The dashboard functionality of Lissio is _not_ the #1 priority. The tool was created to help give users of Kubernetes more information  in an easier fashion than _kubectl get_ or _kubectl describe_
- Lissio can be extended with plugins 
    - Plugin docs here: [docs/plugins](https://github.com/kubenext/lissio/tree/master/docs/plugins)
- Lissio is being very actively developed, with [major releases happening rapidly](https://github.com/kubenext/lissio/releases)
- Lissio is based on newer web technologies. The Kubernetes dashboard is based on "AngularJS" which has been superseded by "Angular". 

## Q: How do I install or Update Lissio?

### Installation:
Lissio can be installed as a package using a variety of package managers, as a pre-built binary, or by building from source.

### Package (Linux only)

1. Download the `.deb` or `.rpm` from the [releases page](https://github.com/kubenext/lissio/releases).

2. Install with either `dpkg -i` or `rpm -i` respectively.

###  Windows

#### Chocolatey

1. Install using chocolatey with the following one-liner:

   ```sh
   choco install lissio --confirm
   ```

#### Scoop

1. Add the [extras](https://github.com/lukesampson/scoop-extras) bucket.

   ```sh
   scoop bucket add extras
   ```

 2. Install using scoop.

   ```sh
   scoop install lissio
   ```

### macOS

#### Homebrew

1. Install using Homebrew with the following one-liner:

   ```sh
   brew install lissio
   ```

### Download a Pre-built Binary (Linux, macOS, Windows)

1. Open the [releases page](https://github.com/kubenext/lissio/releases) from a browser and download the latest tarball or zip file.

2. Extract the tarball or zip where `X.Y` is the release version:

    ```sh
    $ tar -xzvf ~/Downloads/lissio_0.X.Y_Linux-64bit.tar.gz
    lissio_0.X.Y_Linux-64bit/README.md
    lissio_0.X.Y_Linux-64bit/lissio
    ```

3. Verify it runs:

    ```sh
    $ ./lissio_0.X.Y_Linux-64bit/lissio version
    ```

### Building from Source

Lissio can be built from source with the 'Quick Start' instructions found here: [https://github.com/kubenext/lissio/blob/master/HACKING.md](https://github.com/kubenext/lissio/blob/master/HACKING.md)

### Upgrading

The process of upgrading Lissio will depend on how you installed it. Generally, you can use Update or Upgrade functions of the package manager you used to install Lissio. (e.g. brew upgrade lissio)

If you downloaded a pre-built binary, you could download the new version and replace the old one manually.

If you built from source, you would pull the latest from the remote origin (master or the specific release branch), and re-run the *make* build command (e.g. make ci-quick)

## Q: How can I contribute to Lissio?

Lissio is a community-driven project, and as such welcomes new contributors from the community. 

Ways you can contribute with a Pull Request:
- Documentation
    - See something wrong or missing from our docs? 
    - Do you have a unique use-case not documented?
- Lissio core
    - Lissio is written mostly in Golang and Angular. Our hacking guide can be found [here](https://github.com/kubenext/lissio/blob/master/HACKING.md)
- Plugins
    - Lissio has a very extensible plugin model designed to let contributors add functionality. A plugin can read objects, and allows users to add components to Lissio's views.  
    - A sample plugin is available [here](https://github.com/kubenext/lissio/blob/master/cmd/lissio-sample-plugin)
    - A list of community plugins for Lissio will be assembled soon

New contributors will need to sign a CLA (contributor license agreement). We also ask that a changelog entry is included with your pull request. Details are described in our [contributing](https://github.com/kubenext/lissio/blob/master/CONTRIBUTING.md) documentation.

See our [hacking](https://github.com/kubenext/lissio/blob/master/HACKING.md) guide for getting your development environment setup.

See our [roadmap](../ROADMAP.md) for tentative features in a 1.0 release.

**Ways to contribute without a Pull Request**

- Share the love on social media with the hashtag #lissio
- Participate in Lissio community meetings
- Use Lissio and [file issues](https://github.com/kubenext/lissio/issues ) 

## Q: Is Lissio stable?

Lissio is under active development, but each release is considered stable. 

Release information can be found here:
- [Releases](https://github.com/kubenext/lissio/releases)

Open Issues can be found here: 
- [Open Issues](https://github.com/kubenext/lissio/issues)

## Q: Can Lissio connect to multiple clusters at the same time?

No.

Lissio can only connect to a single cluster at a time, but for convenience provides a context switcher that allows you to select the current context without the need to restart.

## Q: Why doesn't Lissio support Feature X?

Lissio is a community driven project with contributions from volunteers around the world. 

If a feature you want is not already on our [Roadmap](https://github.com/kubenext/lissio/blob/master/ROADMAP.md), please feel free to [file an issue](https://github.com/kubenext/lissio/issues/new) and request it, or submit a Pull Request with your feature to be reviewed and [merged](https://github.com/kubenext/lissio/blob/master/CONTRIBUTING.md).

## Q: When will Lissio get Feature X?

See our [roadmap](../ROADMAP.md) for tentative features in a 1.0 release.

## Q: What are the system requirements to run Lissio?

Lissio supports running on macOS, Windows and Linux using either [pre-built binaries](https://github.com/kubenext/lissio/releases) or [building directly from source.](https://github.com/kubenext/lissio/blob/master/HACKING.md)

Lissio requires an active KUBECONFIG (i.e. kubectl configured and working).

Lissio does not reqiure special permissions within your cluster because it uses your local kubeconfig information.

## Q: How can I configure Lissio to run in my Cluster?

While Lissio is designed to run on a developers desktop or laptop, it is possible to configure Lissio to run inside a Kubernetes Cluster. Instructions are located here: https://github.com/kubenext/lissio/tree/master/examples/in-cluster

## Q: Where can I get help with Lissio?

The best way to get help is to file an issue on GitHub:
- https://github.com/kubenext/lissio/issues 

You can also reach out in our communities:

- On Slack 
    - slack.k8s.io #lissio
- On Twitter 
    - [@ProjectLissio](https://twitter.com/projectlissio) 
    - Hashtag:  [#lissio](https://twitter.com/search?q=%23lissio)
- In Google Groups
    - https://groups.google.com/forum/#!forum/project-lissio/

## Q: Where is the Lissio community?

We welcome community engagement in the following places:

- On Slack 
    - slack.k8s.io #lissio
- On Twitter 
    - [@ProjectLissio](https://twitter.com/projectlissio) 
    - Hashtag:  [#lissio](https://twitter.com/search?q=%23lissio)
- In Google Groups
    - https://groups.google.com/forum/#!forum/project-lissio/

