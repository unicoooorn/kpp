# kpp - a lightweight tool for imposing limits on docker volumes usage

kpp (from Russian: КПП, контрольно-пропускной пункт, meaning "border point") is a lightweight tool dedicated to limit docker volume consumption by containers.

It’s especially useful when:
* You cannot use more complex tools like Kubernetes;
* You cannot control or modify the behavior of the containerized programs (e.g., automated testing systems in universities).

Instead of trying to impose limits inside containers, kpp monitors and enforces limits externally.

## Configuration

First, create a configuration file according to [the config structure](https://github.com/unicoooorn/kpp/blob/master/internal/config/config.go).

### Action types

These actions are triggered when a container exceeds its configured volume limit:
* Kill – Immediately terminates the container
* Pause – Pauses the container
* Stop – Gracefully stops the container
* Restart – Restarts the container

### Whitelist & Blacklist
* Whitelist: Add file paths that should be excluded from total volume usage calculation.
* Blacklist: Add file paths that, if detected in a container volume, will cause the container to be immediately killed.

## Usage:
```bash
kpp -c <config_file>
```

## Testing with litterer
To test kpp, you can use build/litterer.py. 
It instantly fills a volume with random files constituting provided number in size
```bash
docker build -t litterer ./build/
```
To run litterer do
```bash
docker run -v <folder_to_litter>:/data litterer <space_amount>
```
