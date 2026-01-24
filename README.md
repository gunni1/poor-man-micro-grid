# Welcome
This is a toy project to learn about technology and practice software developement and DevOps.

# Poor Man's Micro Grid

The Poor Man's Micro Grid is a bunch of software and infrastructure as code components to simulate assets and management software for a small power grid. It operate of assert management level. There is no need for physics simulation and the choosen interfaces and protocols are close to industrial standards but sometimes abstract to focus on learning.

# Components

![Components](doc/micro-grid.drawio.png)

## Assets
Small independend software components which simulat an asset which produce or consume power in the grid. It sends telemetry data and have a interface to receive control commands to increase or decrease suppy/drain or, in case of a battery, switch modes.

## Telemetry
Telemetry components are responsible for collect, store and display telemetry data send by assets. There are stored as timeseries.

## Control
The control plain read telemetry data of the grid to
- evaluate the state of the grid
- make decisions 
- send/write control commands
