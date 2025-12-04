## How to run
Prerequisites: golang and python runtime
1. Clone this repo
```
git clone https://github.com/rei-smz/ewot-app.git && cd ewot-vsensor
```
2. Edit the GraphDB repository endpoint and the description directory in `deploy-temperature.sh` as follow.
```
repository=http://[your-host-name]:7200/repositories/ewot-test/statements
descriptionsDirectory=/absolute/path/to/descriptions/generated/
```
3. Execute `deploy-temperature.sh`.
```
./deploy-temperature.sh
```
