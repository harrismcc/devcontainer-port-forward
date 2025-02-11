# About `devcontainer-port-forward`

Devcontainers are a really cool and useful technology, and while technically the standard is open and can be implemented by anyone, VSCode is still the editor with by far the best support. Any users of other code editors, like Neovim.

There is an official tool [devcontainer/cli](https://github.com/devcontainers/cli) that works pretty well for managing and building devcontainers from the cli, but it lacks the ability to forward opened ports like the native VSCode integration does.

That's where this tool comes in! It's a very simple implementation, that essentially opens up a TCP socket connection locally for each port you specify and forwards/proxies that connection to the devcontainer docker container using [socat](https://linux.die.net/man/1/socat).

## Why Not Just Use Docker's Build-In Port Forwarding?

While you could theoretically use dockers built-in features to forward ports, this actually ends up not working in all cases. In the VSCode implementation, a similar proxying strategy to the one implemented here is used, which results in connections from the host machine appearing to originate from the docker container itself. Many devcontainers are setup with this in mind and so services within them may be configured to only allow connections from the devcontainer.


## Usage

The command is used as follows:

```bash
devcontainer-port-forward -c DOCKER_CONTAINER_ID -p COMMA_SEPARATED_LIST_OF_PORTS
```

So, for example:


```bash
devcontainer-port-forward -c 3b14ca4fb8b53d3401f00bda57d2ea9a746aafd21ee982705d6e53510c1ef780 -p 6007,5432,4001,3000,13000,6080,5901,5907,55432,35432,45432,2048,3002
```
