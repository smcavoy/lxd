---
discourse: lxc:11801
---

(network-forwards)=
# How to configure network forwards

```{note}
Network forwards are available for the {ref}`network-ovn` and the {ref}`network-bridge`.
```

```{youtube} https://www.youtube.com/watch?v=B-Uzo9WldMs
```

Network forwards allow an external IP address (or specific ports on it) to be forwarded to an internal IP address (or specific ports on it) in the network that the forward belongs to.

This feature can be useful if you have limited external IP addresses and want to share a single external address between multiple instances.
There are two different ways how you can use network forwards in this case:

- Forward all traffic from the external address to the internal address of one instance.
  This method makes it easy to move the traffic destined for the external address to another instance by simply reconfiguring the network forward.
- Forward traffic from different port numbers of the external address to different instances (and optionally different ports on those instances).
  This method allows to "share" your external IP address and expose more than one instance at a time.

```{tip}
Network forwards are very similar to using a {ref}`proxy device<devices-proxy>` in NAT mode.

The difference is that network forwards are applied on a network level, while a proxy device is added for an instance.
In addition, proxy devices can be used to proxy traffic between different connection types (for example, TCP and Unix sockets).
```

## List network forwards

View a list of all forwards configured on a network:

`````{tabs}
````{group-tab} CLI

```
lxc network forward list <network_name>
```

Example:

```
lxc network forward list lxdbr0
```

```{note}
This list displays the listen address of the network forward and its default target address, if set. To view the target addresses for a network forward's ports {ref}`set in its port specifications <network-forwards-port-specifications>`, you can {ref}`show details about the network forward <network-forward-show>` or {ref}`edit the network forward <network-forward-edit>`.
```
````

````{group-tab} UI

In {ref}`the web UI <access-ui>`, select {guilabel}`Networks` in the left sidebar, then select the desired network. On the resulting screen, view the {guilabel}`Forwards` tab:

```{figure} /images/UI/forwards_view.png
:width: 95%
:alt: View a list of forwards on a network
```

````
`````

(network-forward-show)=
## Show a network forward

Show details about a specific network forward:

`````{tabs}
````{group-tab} CLI

```
lxc network forward show <network_name> <listen_address>
```

Example:

```
lxc network forward show lxdbr0 192.0.2.1
```
````

````{group-tab} UI

In {ref}`the web UI <access-ui>`, select {guilabel}`Networks` in the left sidebar, then select the desired network. On the resulting screen, view the {guilabel}`Forwards` tab. This tab shows you information about all forwards on the network. You can click the {guilabel}`Edit` icon to view details for a specific forward:

```{figure} /images/UI/forward_edit_ex1.png
:width: 95%
:alt: View details about a specific forward on a network through its edit screen
```

````
`````

(network-forward-create)=
## Create a network forward

(network-forwards-listen-addresses)=
### Requirements for listen addresses

Before you can create a network forward, you must understand the requirements for listen addresses.

For both OVN and bridge networks, the listen addresses must not overlap with any subnet in use by other networks on the host. Otherwise, the listen address requirements differ by network type.

````{tabs}

```{group-tab} OVN network

For an OVN network, the allowed listen addresses must be defined in at least one of the following configuration options, using [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing):

- {config:option}`network-bridge-network-conf:ipv4.routes` or {config:option}`network-bridge-network-conf:ipv6.routes` in the OVN network's uplink network configuration
- {config:option}`project-restricted:restricted.networks.subnets` in the OVN network's project configuration

```

```{group-tab} Bridge network

A bridge network does not require you to define allowed listen addresses. Use any non-conflicting IP address available on the host.

```

````

### Create a forward in an OVN network

```{note}
You must configure the {ref}`allowed listen addresses <network-forwards-listen-addresses>` before you can create a forward in an OVN network.

The IP addresses and ports shown in the examples below are only examples. It is up to you to choose the allowed and available addresses and ports for your setup.
```

`````{tabs}
````{group-tab} CLI

Use the following command to create a forward in an OVN network:

```
lxc network forward create <ovn_network_name> <listen_address>|--allocate=ipv{4,6} [target_address=<target_address>] [user.<key>=<value>]
```

- For `<ovn_network_name>`, specify the name of the OVN network on which to create the forward.
- Immediately following the network name, provide only one of the following for the listen address:
   - A listen IP address allowed by the {ref}`network-forwards-listen-addresses` (no port number)
   - The `--allocate=` flag with a value of either `ipv4` or `ipv6` for automatic allocation of an allowed IP address
- Optionally provide a default `target_address` (no port number). Any traffic that does not match a port specification is forwarded to this address. This must be an IP address within the OVN network's subnet; typically, the static IP address of an instance is used.
- Optionally provide custom user.* keys to be stored in the network forward's configuration.

This example shows how to create a network forward on a network named `ovn1` with an allocated listen address and no default target address:

```
lxc network forward create ovn1 --allocate=ipv4
```

This example shows how to create a network forward on a network named `ovn1` with a specific listen address and a target address:

```
lxc network forward create ovn1 192.0.2.1 target_address=10.41.211.2
```

````

````{group-tab} UI

In {ref}`the web UI <access-ui>`, select {guilabel}`Networks` in the left sidebar, then select the desired OVN network. On the resulting screen, view the {guilabel}`Forwards` tab. Click the {guilabel}`Create forward` button.

In the {guilabel}`Create a new forward` panel, only the {guilabel}`Listen address` field is required.

```{figure} /images/UI/forward_create_ovn.png
:width: 95%
:alt: Create an OVN network forward
```

- For the {guilabel}`Listen address`, provide an IP address allowed by the {ref}`network-forwards-listen-addresses` (no port number).
- Optionally provide a {guilabel}`Default target address` (no port number). Any traffic that does not match a port specification is forwarded to this address. This must be an IP address within the OVN network's subnet; typically, the static IP address of an instance is used.

You can optionally set up port specifications for the network forward by clicking the {guilabel}`Add port` button. These specifications allow forwarding traffic from specific ports on the listen address to ports on a target address. For details on how to configure this section, see: {ref}`network-forwards-port-specifications`.

Once you have finished setting up the network forward, click the {guilabel}`Create` button.

````
`````

### Create a forward in a bridge network

```{note}
The IP addresses and ports shown in the examples below are only examples. It is up to you to choose the allowed and available addresses and ports for your setup.
```

`````{tabs}
````{group-tab} CLI

Use the following command to create a forward in a bridge network:

```
lxc network forward create <bridge_network_name> <listen_address> [target_address=<target_address>] [user.<key>=<value>]
```

- For `<bridge_network_name>`, specify the name of the bridge network on which to create the forward.
- Immediately following the network name, provide an IP address allowed by the {ref}`network-forwards-listen-addresses` (no port number).
- Optionally provide a default `target_address` (no port number). Any traffic that does not match a port specification is forwarded to this address. This must be an IP address within the bridge network's subnet; typically, the static IP address of an instance is used.
- Optionally provide custom user.* keys to be stored in the network forward's configuration.
- You cannot use the `--allocate` flag with bridge networks.

This example shows how to create a forward on a network named `bridge1` with a specific listen address and a target address:

```
lxc network forward create bridge1 192.0.2.1 target_address=10.41.211.2
```

````

````{group-tab} UI

In {ref}`the web UI <access-ui>`, select {guilabel}`Networks` in the left sidebar, then select the desired bridge network. On the resulting screen, view the {guilabel}`Forwards` tab. Click the {guilabel}`Create forward` button.

In the {guilabel}`Create a new forward` panel, only the {guilabel}`Listen address` field is required.

```{figure} /images/UI/forward_create_bridge.png
:width: 95%
:alt: Create a bridge network forward
```

- For the {guilabel}`Listen address`, provide a listen IP address allowed by the {ref}`network-forwards-listen-addresses` (no port number).
- Optionally provide a {guilabel}`Default target address` (no port number). Any traffic that does not match a port specification is forwarded to this address. This must be an IP address within the bridge network's subnet; typically, the static IP address of an instance is used.

You can optionally set up port specifications for the network forward by clicking the {guilabel}`Add port` button. These specifications allow forwarding traffic from specific ports on the listen address to ports on a target address. For details on how to configure this section, see: {ref}`network-forwards-port-specifications`.

Once you have finished setting up the network forward, click the {guilabel}`Create` button.

````
`````

### Forward properties

Network forwards have the following properties:

% Include content from [../metadata.txt](../metadata.txt)
```{include} ../metadata.txt
    :start-after: <!-- config group network-forward-forward-properties start -->
    :end-before: <!-- config group network-forward-forward-properties end -->
```

(network-forwards-port-specifications)=
## Configure ports

Once a forward is created on a network (whether bridge or OVN), it can be configured with port specifications. These specifications allow forwarding traffic from ports on the listen address to ports on a target address.

`````{tabs}
````{group-tab} CLI

When using the CLI, you must first {ref}`create a network forward <network-forward-create>` before you can add port specifications to it.

Use the following command to add port specifications on a forward:

```
lxc network forward port add <network_name> <listen_address> <protocol> <listen_ports> <target_address> [<target_ports>]
```

- Use the network name and listen address of the forward for which you want to add port specifications.
- Use either `tcp` or `udp` as the protocol.
- For the listen ports, you can specify a single listen port, a port range, or a comma-separated set of ports/port ranges.
- Specify a target address. This address must be within the network's subnet, and it must be different from the forward's default target address. Typically, the static IP address of an instance is used.
- Optionally specify a target port or ports. You can:
   - Specify a single target port to forward traffic from all listen ports to this target port.
   - Specify a set of target ports with the same number of set items as the listen ports. This forwards traffic from the first listen port to the first target port, the second listen port to the second target port, and so on.
- If no target port is specified, the listen port value is used for the target port.

The example below shows how to configure a forward with a single listen port. Since no target port is specified, the target port defaults to the value of the listen port:

```
lxc network forward port add network1 192.0.2.1 tcp 22 10.41.211.2
```

The example below shows how to configure a forward with a set of listen ports mapped to a single target port. Traffic to the listen address at ports 80 and 90 through 100 is forwarded to port 443 of the target address:

```
lxc network forward port add network1 192.0.2.1 tcp 80,90-100 10.41.211.2 443
```

The example below shows how to configure a forward with a set of listen ports mapped to a set of target ports. Traffic to the listen address at port 22 is forwarded to port 22 of the target address, whereas traffic to port 80 is forwarded to port 443:

```
lxc network forward port add network1 192.0.2.1 tcp 22,80 10.41.211.2 22,443
```

````

````{group-tab} UI

These specifications allow forwarding traffic from ports on the listen address to ports on a target address. In the web UI, you can configure port specifications on a network forward at the time you {ref}`create the forward <network-forward-create>`, or by {ref}`editing the forward <network-forward-edit>` after creation.

```{figure} /images/UI/forward_create_port.png
:width: 95%
:alt: Configure a network forward's port specifications
```

- For the {guilabel}`Listen port`, you can specify a single port, a port range, or a comma-separated set of ports/port ranges.
- Select either {guilabel}`TCP` or {guilabel}`UDP` as the protocol.
- Specify a {guilabel}`Target address`. This address must be within the network's subnet, and it must be different from the forward's {guilabel}`Default target address`. Typically, the static IP address of an instance is used.
- Optionally specify a target port or ports. You can:
   - Specify a single target port to forward traffic from all listen ports to this target port.
   - Specify a set of target ports with the same number of set items as the listen ports. This forwards traffic from the first listen port to the first target port, the second listen port to the second target port, and so on.
- If no target port is specified, the listen port value is used for the target port.

Examples:

- If the {guilabel}`Listen port` is set to `22` and no {guilabel}`Target port` is specified, the target port value defaults to 22.
- If the {guilabel}`Listen port` is set to `80,90-100` and the {guilabel}`Target port` is set to {guilabel}`442`, all traffic to the listen address at ports 80 and 90 through 100 is forwarded to port 443 of the target address.
- If the {guilabel}`Listen port` is set to `22,80` and the {guilabel}`Target port` is set to `22,443`, all traffic to the listen address at port 22 is forwarded to port 22 of the target address, whereas traffic to port 80 is forwarded to port 443.

````
`````

### Port properties

Network forward ports have the following properties:

% Include content from [../metadata.txt](../metadata.txt)
```{include} ../metadata.txt
    :start-after: <!-- config group network-forward-port-properties start -->
    :end-before: <!-- config group network-forward-port-properties end -->
```

(network-forward-edit)=
## Edit a network forward

`````{tabs}
````{group-tab} CLI

Use the following command to edit a network forward:

```bash
lxc network forward edit <network_name> <listen_address>
```

This command opens the network forward in YAML format for editing.
You can edit both the general configuration and the port specifications.

````

````{group-tab} UI

In {ref}`the web UI <access-ui>`, select {guilabel}`Networks` in the left sidebar, then select the desired network. On the resulting screen, view the {guilabel}`Forwards` tab. This tab shows you information about all forwards on the network. Click the {guilabel}`Edit` icon next to a forward to edit it:

```{figure} /images/UI/forward_edit_ex1.png
:width: 95%
:alt: Choose to edit a forward on a network
```

In the resulting screen, you can edit the forward's general configuration as well as its port specifications:

```{figure} /images/UI/forward_edit_ex2.png
:width: 95%
:alt: Edit a forward on a network
```

````
`````

## Delete a network forward

`````{tabs}
````{group-tab} CLI

Use the following command to delete a network forward:

```bash
lxc network forward delete <network_name> <listen_address>
```

````

````{group-tab} UI

In {ref}`the web UI <access-ui>`, select {guilabel}`Networks` in the left sidebar, then select the desired network. On the resulting screen, view the {guilabel}`Forwards` tab. This tab shows you information about all forwards on the network. Click the {guilabel}`Delete` icon next to a forward to delete it:

```{figure} /images/UI/forward_delete.png
:width: 95%
:alt: Delete a forward on a network
```

````
`````
