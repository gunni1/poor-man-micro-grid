# Introduction
This piece of code is about simulating a PV asset within a micro grid demo environment. It uses MQTT to export its telemetry data to

```
microgrid/<ASSET_TYPE>/<ASSET_ID>/telemetry
```

# Usage

## Configuration
Following environment variables can be used to configre the pv simulation:


- MQTT_BROKER: Adress of the MQTT broker to send the simulated telemetry data
- ASSET_ID: Name of the asset to identify it
- P_NOMINAL_KW: Nominal power of the pv